package dto

import (
	"time"

	"github.com/google/uuid"
)

type TrackEpisodeResponse struct {
	TrackID         uuid.UUID `json:"track_id"`
	Artist          string    `json:"artist"`
	Title           string    `json:"title"`
	CurrentPosition int       `json:"current_position"`
}

type EpisodeResponse struct {
	ID        uuid.UUID              `json:"id"`
	ChartID   uuid.UUID              `json:"chart_id"`
	CreatedAt time.Time              `json:"created_at"`
	Tracks    []TrackEpisodeResponse `json:"tracks"`
}
