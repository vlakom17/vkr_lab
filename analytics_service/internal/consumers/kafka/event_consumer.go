package kafka

import (
	"context"
	"encoding/json"
	"log"

	"charts-analytics-service/internal/domain/event"
	"charts-analytics-service/internal/infrastructure/mb"
)

type KafkaConsumer struct {
	consumer *mb.Consumer
}

func NewKafkaConsumer(c *mb.Consumer) *KafkaConsumer {
	return &KafkaConsumer{consumer: c}
}

func (c *KafkaConsumer) ConsumeReactions(
	ctx context.Context,
	handler func(context.Context, event.ReactionEvent) error,
) {
	c.consumer.Start(ctx, func(data []byte) error {

		var r event.ReactionEvent

		if err := json.Unmarshal(data, &r); err != nil {
			log.Println("failed to unmarshal reaction event:", err)
			return err
		}

		if err := handler(ctx, r); err != nil {
			log.Println("handler error:", err)
			return err
		}

		return nil
	})
}
