package dto

import (
	"time"

	"github.com/google/uuid"
)

type TrackEpisodeResponse struct {
	TrackID             uuid.UUID `json:"track_id"`
	Artist              string    `json:"artist"`
	Title               string    `json:"title"`
	CurrentPosition     int       `json:"current_position"`
	PreviousPosition    int       `json:"previous_position"`
	HighestPosition     int       `json:"highest_position"`
	EpisodesCount       int       `json:"episodes_count"`
	TimesAtPeakPosition int       `json:"times_at_peak_position"`
}

type EpisodeResponse struct {
	ID        uuid.UUID              `json:"id"`
	ChartID   uuid.UUID              `json:"chart_id"`
	CreatedAt time.Time              `json:"created_at"`
	Tracks    []TrackEpisodeResponse `json:"tracks"`
}
