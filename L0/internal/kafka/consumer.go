package kafka

import (
	"context"
	"encoding/json"
	"log"
	"myapp/internal/cache"
	"myapp/internal/model"
	"myapp/internal/repository"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	repo   *repository.OrderRepository
	cache  *cache.OrderCache
}

func NewConsumer(brokers []string, topic string, repo *repository.OrderRepository, cache *cache.OrderCache) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "order-group",
	})
	return &Consumer{reader: r, repo: repo, cache: cache}
}

func (c *Consumer) Run(ctx context.Context) {
	log.Println("Kafka consumer started...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping...")
			c.reader.Close()
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from Kafka: %v", err)
				continue
			}

			var order model.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				log.Printf("Error unmarshalling message: %v. Message: %s", err, string(m.Value))

				continue
			}

			if order.OrderUID == "" {
				log.Println("Invalid order received (empty order_uid), skipping.")
				continue
			}

			order.Delivery.OrderUID = order.OrderUID
			order.Payment.Transaction = order.OrderUID
			for i := range order.Items {
				order.Items[i].OrderUID = order.OrderUID
			}

			if err := c.repo.SaveOrder(ctx, order); err != nil {
				log.Printf("Error saving order to DB: %v", err)

				continue
			}

			c.cache.Set(order)
			log.Printf("Successfully processed and cached order: %s", order.OrderUID)
		}
	}
}
