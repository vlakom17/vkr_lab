package kafka

import (
	"context"
	"encoding/json"
	"log"

	"charts-archive-service/internal/domain/event"
	"charts-archive-service/internal/infrastructure/mb"
)

type KafkaConsumer struct {
	consumer *mb.Consumer
}

func NewKafkaConsumer(c *mb.Consumer) *KafkaConsumer {
	return &KafkaConsumer{consumer: c}
}

func (c *KafkaConsumer) ConsumeEpisodes(
	ctx context.Context,
	handler func(context.Context, event.EpisodeSnapshotEvent) error,
) {
	c.consumer.Start(ctx, func(data []byte) error {

		var e event.EpisodeSnapshotEvent

		if err := json.Unmarshal(data, &e); err != nil {
			log.Println("failed to unmarshal episode event:", err)
			return err
		}

		if err := handler(ctx, e); err != nil {
			log.Println("handler error:", err)
			return err
		}

		return nil
	})
}
