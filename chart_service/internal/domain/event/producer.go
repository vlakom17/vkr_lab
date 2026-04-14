package event

import "context"

type EventProducer interface {
	SendReaction(ctx context.Context, event ReactionEvent) error
	SendEpisode(ctx context.Context, event EpisodeSnapshotEvent) error
}
