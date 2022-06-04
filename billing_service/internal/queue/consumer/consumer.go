package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/service"
	invoicePaidConsumer "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/queue/consumer/invoice_paid"
	invoicePayFailedConsumer "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/queue/consumer/invoice_pay_failed"
	reservationCreatedConsumer "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/queue/consumer/reservation_created"
)

const (
	groupID = "billing_service"
)

func Consume(ctx context.Context, brokers []string, s service.Invoices) error {
	cfg := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return fmt.Errorf("create producer: %w", err)
	}

	go func() {
		select {
		case <-producer.Errors():
			fmt.Printf("producer err: %s", err)
		case <-ctx.Done():
			return
		}
	}()

	cg, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	handlers := map[string]TopicHandler{
		// Event of successful reservation of order products
		"reservation_created": reservationCreatedConsumer.New(producer, s),

		// Invoice payment successful event (imitation of user action)
		"invoice_paid": invoicePaidConsumer.New(producer, s),
		// Invoice payment failed event (imitation of user action)
		"invoice_pay_failed": invoicePayFailedConsumer.New(producer, s),
	}

	registerHandler(ctx, cg, handlers)

	return nil
}

func registerHandler(ctx context.Context, cg sarama.ConsumerGroup, handlers map[string]TopicHandler) {
	var topics []string
	for topic := range handlers {
		topics = append(topics, topic)
	}

	go func() {
		for {
			log.Printf("Consume %v topic...", topics)
			err := cg.Consume(ctx, topics, NewHandler(handlers))
			if err != nil {
				log.Printf("consume error: %s", err)
				time.Sleep(time.Second * 5)
			}
		}
	}()
}
