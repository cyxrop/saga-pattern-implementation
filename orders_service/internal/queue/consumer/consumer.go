package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/service"
	createOrderConsumer "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/queue/consumer/create_order"
	orderPaidConsumer "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/queue/consumer/order_paid"
	orderPayFailedConsumer "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/queue/consumer/order_pay_failed"
	reservationCreatedConsumer "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/queue/consumer/reservation_created"
	reservationFailedConsumer "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/queue/consumer/reservation_failed"
)

const (
	groupID = "orders_service"
)

func Consume(ctx context.Context, brokers []string, s service.Orders) error {
	log.Println("brokers: ", brokers)

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
		// Requests to create orders
		"create_order": createOrderConsumer.New(producer, s),

		// Event of successful reservation of order products
		"reservation_created": reservationCreatedConsumer.New(producer, s),
		// Event of failed reservation of order products
		"reservation_failed": reservationFailedConsumer.New(producer, s),

		// Successful order payment event
		"order_paid": orderPaidConsumer.New(producer, s),
		// Failed order payment event
		"order_pay_failed": orderPayFailedConsumer.New(producer, s),
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
