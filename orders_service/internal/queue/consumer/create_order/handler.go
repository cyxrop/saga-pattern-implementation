package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/service"
)

type CreateOrderHandler struct {
	producer sarama.AsyncProducer
	os       service.Orders
}

func New(producer sarama.AsyncProducer, os service.Orders) *CreateOrderHandler {
	return &CreateOrderHandler{producer: producer, os: os}
}

func (h *CreateOrderHandler) Handle(msg *sarama.ConsumerMessage) {
	var m CreateOrderMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	products := make([]model.OrderProduct, len(m.Products))
	for i, pr := range m.Products {
		products[i] = model.OrderProduct{
			ProductID: pr.ProductID,
			Number:    pr.Number,
		}
	}

	ID, err := h.os.CreateOrder(
		context.Background(),
		model.Order{WarehouseID: m.WarehouseID},
		products,
	)
	if err != nil {
		log.Printf("hande order create failed: %s", err)
		return
	}

	log.Printf("order created: %v", ID)

	orderCreatedMsg := OrderCreatedMessage{
		OrderID:     ID,
		WarehouseID: m.WarehouseID,
		Products:    m.Products,
	}

	marshaled, err := json.Marshal(orderCreatedMsg)
	if err != nil {
		log.Printf("failed to marshal struct %v: %s", orderCreatedMsg, err)
		return
	}

	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: "order_created",
		Key:   sarama.StringEncoder(strconv.FormatInt(ID, 10)),
		Value: sarama.ByteEncoder(marshaled),
	}
	log.Printf("produced message to order_created: %v", ID)
}
