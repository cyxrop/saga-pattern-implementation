package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/service"
)

type ReservationCreatedHandler struct {
	producer sarama.AsyncProducer
	is       service.Invoices
}

func New(producer sarama.AsyncProducer, is service.Invoices) *ReservationCreatedHandler {
	return &ReservationCreatedHandler{producer: producer, is: is}
}

func (h *ReservationCreatedHandler) Handle(msg *sarama.ConsumerMessage) {
	var m ReservationCreatedMessage
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		log.Printf("unmarshal message %s: %s", string(msg.Value), err)
		return
	}

	ID, err := h.is.CreateInvoice(context.Background(), m.OrderID, m.Amount)
	if err != nil {
		log.Printf("failed to create order %d invoice: %s", m.OrderID, err)
		h.produceMessage("order_pay_failed", m.OrderID, msg.Value)
		return
	}

	log.Printf("invoice %d of order %d issued", ID, m.OrderID)

	issuedMsg := InvoiceIssuedMessage{
		InvoiceID: ID,
		OrderID:   m.OrderID,
		Amount:    m.Amount,
	}

	marshaled, err := json.Marshal(issuedMsg)
	if err != nil {
		log.Printf("marshal invoice %d issued message failed: %s", ID, err)
		return
	}

	h.produceMessage("invoice_issued", m.OrderID, marshaled)
	log.Printf("produced message to invoice_issued, order %d", m.OrderID)
}

func (h *ReservationCreatedHandler) produceMessage(topic string, key int64, value []byte) {
	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(key, 10)),
		Value: sarama.ByteEncoder(value),
	}
}
