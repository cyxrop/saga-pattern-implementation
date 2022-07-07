package consumer

import (
	"github.com/Shopify/sarama"
)

type TopicHandler interface {
	Handle(*sarama.ConsumerMessage)
}

type Handler struct {
	handlers map[string]TopicHandler
}

func NewHandler(handlers map[string]TopicHandler) *Handler {
	return &Handler{
		handlers: handlers,
	}
}

func (h *Handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h, ok := h.handlers[msg.Topic]
		if !ok {
			continue
		}

		h.Handle(msg)
	}

	return nil
}
