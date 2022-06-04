package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/service"
	orderCreatedConsumer "gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/queue/consumer/order_created"
	orderPayFailedConsumer "gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/queue/consumer/order_pay_failed"
)

const (
	groupID = "products_service"
)

func Consume(ctx context.Context, brokers []string, s service.Products) error {
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
		// Event of successful order creation
		"order_created": orderCreatedConsumer.New(producer, s),
		// Failed order payment event
		"order_pay_failed": orderPayFailedConsumer.New(s),
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
