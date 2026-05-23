package response

import (
	"time"

	"github.com/google/uuid"
)

type ListenLinks struct {
	AppleMusic  string `json:"apple_music"`
	YandexMusic string `json:"yandex_music"`
}

type TrackEpisodeResponse struct {
	TrackID             uuid.UUID   `json:"track_id"`
	Artist              string      `json:"artist"`
	Title               string      `json:"title"`
	CurrentPosition     int         `json:"current_position"`
	PreviousPosition    int         `json:"previous_position"`
	HighestPosition     int         `json:"highest_position"`
	TimesAtPeakPosition int         `json:"times_at_peak_position"`
	EpisodesCount       int         `json:"episodes_count"`
	ListenLinks         ListenLinks `json:"listen_links"`
}

type EpisodeResponse struct {
	ID        uuid.UUID              `json:"id"`
	ChartID   uuid.UUID              `json:"chart_id"`
	CreatedAt time.Time              `json:"created_at"`
	Tracks    []TrackEpisodeResponse `json:"tracks"`
}
