package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/service"
)

type OrderPaidHandler struct {
	producer sarama.AsyncProducer
	os       service.Orders
}

func New(producer sarama.AsyncProducer, os service.Orders) *OrderPaidHandler {
	return &OrderPaidHandler{producer: producer, os: os}
}

func (h *OrderPaidHandler) Handle(msg *sarama.ConsumerMessage) {
	var m OrderPaidMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.os.UpdateOrderStatus(context.Background(), m.OrderID, model.OrderStatusPaid, "")
	if err != nil {
		log.Printf("failed to update order %d status to %d: %s", m.OrderID, model.OrderStatusPaid, err)
		return
	}

	log.Printf("order %d status has been updated to %d", m.OrderID, model.OrderStatusPaid)
}
