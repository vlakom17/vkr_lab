package event

import (
	"time"

	"github.com/google/uuid"
)

type ReactionType string

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
	ReactionView    ReactionType = "view"
	ReactionRemove  ReactionType = "remove"
)

type ReactionEvent struct {
	UserID    uuid.UUID    `json:"user_id"`
	ChartID   uuid.UUID    `json:"chart_id"`
	Type      ReactionType `json:"type"`
	CreatedAt time.Time    `json:"created_at"`
}
