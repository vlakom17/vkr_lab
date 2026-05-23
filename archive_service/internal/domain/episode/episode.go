package episode

import (
	"time"

	"charts-archive-service/internal/domain/track_episode"

	"github.com/google/uuid"
)

type Episode struct {
	ID        uuid.UUID
	ChartID   uuid.UUID
	CreatedAt time.Time
	Tracks    []track_episode.TrackEpisode
}
