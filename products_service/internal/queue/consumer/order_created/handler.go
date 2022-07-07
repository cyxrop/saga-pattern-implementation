package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/service"
)

type OrderCreatedHandler struct {
	producer sarama.AsyncProducer
	ps       service.Products
}

func New(producer sarama.AsyncProducer, ps service.Products) *OrderCreatedHandler {
	return &OrderCreatedHandler{producer: producer, ps: ps}
}

func (h *OrderCreatedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m OrderCreatedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	productItems := make([]model.ProductItem, len(m.Products))
	for i, pr := range m.Products {
		productItems[i] = model.ProductItem{
			ProductID: pr.ProductID,
			Number:    pr.Number,
		}
	}

	ctx := context.Background()
	orderAmount, err := h.ps.CalculateAmount(ctx, productItems)
	if err != nil {
		log.Printf("reservation of %d order failed at the order amount calucation stage: %s", m.OrderID, err)
		h.failReservation(m.OrderID, msg.Value)
		return
	}

	reservationsIDs, err := h.ps.ReserveOrderProducts(ctx, m.OrderID, m.WarehouseID, productItems)
	if err != nil {
		log.Printf("reservation of %d order failed at the products reservation stage: %s", m.OrderID, err)
		h.failReservation(m.OrderID, msg.Value)
		return
	}

	log.Printf("created reservations: %v", reservationsIDs)
	h.finishReservation(ReservationCreatedMessage{
		OrderID:        m.OrderID,
		ReservationIDs: reservationsIDs,
		Amount:         orderAmount,
	})
}

func (h *OrderCreatedHandler) failReservation(key int64, value []byte) {
	h.produceMessage("reservation_failed", key, value)
}

func (h *OrderCreatedHandler) finishReservation(msg ReservationCreatedMessage) {
	value, err := json.Marshal(msg)
	if err != nil {
		log.Printf("finish reservation: marshal order created message: %s", err)
		return
	}

	h.produceMessage("reservation_created", msg.OrderID, value)
	log.Printf("produced message to reservation_created: %v", msg.OrderID)
}

func (h *OrderCreatedHandler) produceMessage(topic string, key int64, value []byte) {
	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(key, 10)),
		Value: sarama.ByteEncoder(value),
	}
}
