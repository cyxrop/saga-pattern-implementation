package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/service"
)

type ReservationFailedHandler struct {
	producer sarama.AsyncProducer
	os       service.Orders
}

func New(producer sarama.AsyncProducer, os service.Orders) *ReservationFailedHandler {
	return &ReservationFailedHandler{producer: producer, os: os}
}

func (h *ReservationFailedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m ReservationFailedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.os.UpdateOrderStatus(context.Background(), m.OrderID, model.OrderStatusFailed, "reservation of order failed")
	if err != nil {
		log.Printf("failed to fail order %d: %s", m.OrderID, err)
		return
	}

	log.Printf("order %d failed after reservation fail", m.OrderID)
}
