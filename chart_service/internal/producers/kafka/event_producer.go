package kafka

import (
	"context"
	"encoding/json"

	"charts-chart-service/internal/domain/event"
	"charts-chart-service/internal/infrastructure/mb"
)

type KafkaProducer struct {
	kafka *mb.Kafka
}

func NewKafkaProducer(k *mb.Kafka) *KafkaProducer {
	return &KafkaProducer{kafka: k}
}

func (p *KafkaProducer) SendReaction(ctx context.Context, e event.ReactionEvent) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return p.kafka.SendReaction(ctx, bytes)
}

func (p *KafkaProducer) SendEpisode(ctx context.Context, e event.EpisodeSnapshotEvent) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return p.kafka.SendEpisode(ctx, bytes)
}
