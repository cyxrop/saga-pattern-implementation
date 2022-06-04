package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/service"
)

type ReservationCreatedHandler struct {
	producer sarama.AsyncProducer
	os       service.Orders
}

func New(producer sarama.AsyncProducer, os service.Orders) *ReservationCreatedHandler {
	return &ReservationCreatedHandler{producer: producer, os: os}
}

func (h *ReservationCreatedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m ReservationCreatedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.os.UpdateOrderStatus(context.Background(), m.OrderID, model.OrderStatusPendingPayment, "")
	if err != nil {
		log.Printf("failed to update order %d status to %d: %s", m.OrderID, model.OrderStatusPendingPayment, err)
		return
	}

	log.Printf("order %d status has been updated to %d", m.OrderID, model.OrderStatusPendingPayment)
}
