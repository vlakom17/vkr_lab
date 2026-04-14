package event

import (
	"time"

	"github.com/google/uuid"
)

type ReactionEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	ChartID   uuid.UUID `json:"chart_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
