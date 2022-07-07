package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/service"
)

type OrderPayFailedHandler struct {
	producer sarama.AsyncProducer
	os       service.Orders
}

func New(producer sarama.AsyncProducer, os service.Orders) *OrderPayFailedHandler {
	return &OrderPayFailedHandler{producer: producer, os: os}
}

func (h *OrderPayFailedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m OrderPayFailedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.os.UpdateOrderStatus(context.Background(), m.OrderID, model.OrderStatusFailed, "order payment failed")
	if err != nil {
		log.Printf("failed to update order %d status to %d: %s", m.OrderID, model.OrderStatusFailed, err)
		return
	}

	log.Printf("order %d status has been updated to %d after failed payment", m.OrderID, model.OrderStatusFailed)
}
