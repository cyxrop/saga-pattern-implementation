package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/service"
)

type OrderPayFailedHandler struct {
	producer sarama.AsyncProducer
	os       service.Products
}

func New(ps service.Products) *OrderPayFailedHandler {
	return &OrderPayFailedHandler{os: ps}
}

func (h *OrderPayFailedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m OrderPayFailedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.os.CancelOrderReservation(context.Background(), m.OrderID)
	if err != nil {
		log.Printf("failed to cancel order %d reservations: %s", m.OrderID, err)
		return
	}

	log.Printf("order %d reservations have been canceled after failed order payment", m.OrderID)
}
