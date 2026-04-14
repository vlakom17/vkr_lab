package mb

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	ReactionWriter *kafka.Writer
	EpisodeWriter  *kafka.Writer
}

func NewKafka(broker string) *Kafka {

	reactionWriter := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "reactions",
		Balancer: &kafka.LeastBytes{},
	}

	episodeWriter := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "episodes",
		Balancer: &kafka.LeastBytes{},
	}

	log.Println("connected to kafka")

	return &Kafka{
		ReactionWriter: reactionWriter,
		EpisodeWriter:  episodeWriter,
	}
}

func (k *Kafka) SendReaction(ctx context.Context, msg []byte) error {
	return k.ReactionWriter.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
}

func (k *Kafka) SendEpisode(ctx context.Context, msg []byte) error {
	return k.EpisodeWriter.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
}

func (k *Kafka) Close() {
	if err := k.ReactionWriter.Close(); err != nil {
		log.Println("error closing reaction writer:", err)
	}
	if err := k.EpisodeWriter.Close(); err != nil {
		log.Println("error closing episode writer:", err)
	}
}
