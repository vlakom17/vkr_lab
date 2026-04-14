package event

import (
	"time"

	"github.com/google/uuid"
)

type TrackSnapshot struct {
	Artist          string `json:"artist"`
	Title           string `json:"title"`
	CurrentPosition int    `json:"current_position"`
}

type EpisodeSnapshotEvent struct {
	ChartID   uuid.UUID       `json:"chart_id"`
	Tracks    []TrackSnapshot `json:"tracks"`
	CreatedAt time.Time       `json:"created_at"`
}
