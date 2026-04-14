package mb

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var KafkaTopic = "reactions"

type Consumer struct {
	reader *kafka.Reader
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	log.Println("kafka consumer connected")

	return &Consumer{reader: reader}
}

func (c *Consumer) Start(ctx context.Context, handler func([]byte) error) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("consumer stopped")
				return
			}
			log.Println("error reading message:", err)
			continue
		}

		for i := 0; i < 3; i++ {
			err = handler(msg.Value)
			if err == nil {
				break
			}
			log.Printf("handler retry %d failed: %v", i+1, err)
		}

		if err != nil {
			log.Println("skip message after retries")
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("commit error:", err)
		}
	}
}
