package request

import "charts-chart-service/internal/domain/event"

type CreateEpisodeRequest struct {
	Tracks []event.TrackEntryData `json:"tracks" binding:"required"`
}
