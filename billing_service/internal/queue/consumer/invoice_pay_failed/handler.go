package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/service"
)

type InvoicePayFailedHandler struct {
	producer sarama.AsyncProducer
	is       service.Invoices
}

func New(producer sarama.AsyncProducer, is service.Invoices) *InvoicePayFailedHandler {
	return &InvoicePayFailedHandler{producer: producer, is: is}
}

func (h *InvoicePayFailedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m InvoicePayFailedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.is.UpdateInvoiceStatus(context.Background(), m.InvoiceID, model.InvoiceStatusFailed)
	if err != nil {
		log.Printf("failed to update invoice %d status: %s", m.InvoiceID, err)
		return
	}

	marshaled, err := json.Marshal(OrderPayFailedMessage{
		OrderID: m.OrderID,
	})
	if err != nil {
		log.Printf("marshal order %d pay failed message failed: %s", m.OrderID, err)
		return
	}

	h.produceMessage("order_pay_failed", m.OrderID, marshaled)
	log.Printf("produced message to order_pay_failed, order %d", m.OrderID)
}

func (h *InvoicePayFailedHandler) produceMessage(topic string, key int64, value []byte) {
	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(key, 10)),
		Value: sarama.ByteEncoder(value),
	}
}
