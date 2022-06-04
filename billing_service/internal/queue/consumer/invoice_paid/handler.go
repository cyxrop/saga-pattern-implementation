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

type InvoicePaidHandler struct {
	producer sarama.AsyncProducer
	is       service.Invoices
}

func New(producer sarama.AsyncProducer, is service.Invoices) *InvoicePaidHandler {
	return &InvoicePaidHandler{producer: producer, is: is}
}

func (h *InvoicePaidHandler) Handle(msg *sarama.ConsumerMessage) {
	var m InvoicePaidMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	err := h.is.UpdateInvoiceStatus(context.Background(), m.InvoiceID, model.InvoiceStatusPaid)
	if err != nil {
		log.Printf("failed to update invoice %d status: %s", m.InvoiceID, err)
		h.produceMessage("order_pay_failed", m.OrderID, msg.Value)
		return
	}

	log.Printf("order %d invoice %d paid", m.OrderID, m.InvoiceID)
	orderPaidMsg := OrderPaidMessage{
		InvoiceID: m.InvoiceID,
		OrderID:   m.OrderID,
		Amount:    m.Amount,
	}

	marshaled, err := json.Marshal(orderPaidMsg)
	if err != nil {
		log.Printf("marshal order %d paid message failed: %s", m.OrderID, err)
		return
	}

	h.produceMessage("order_paid", m.OrderID, marshaled)
	log.Printf("produced message to order_paid, order %d", m.OrderID)
}

func (h *InvoicePaidHandler) produceMessage(topic string, key int64, value []byte) {
	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(key, 10)),
		Value: sarama.ByteEncoder(value),
	}
}
