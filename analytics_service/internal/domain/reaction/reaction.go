package reaction

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

type Reaction struct {
	UserID    uuid.UUID
	ChartID   uuid.UUID
	Type      ReactionType
	CreatedAt time.Time
}

type ReactionStats struct {
	ChartID       uuid.UUID
	LikesCount    int
	DislikesCount int
	ViewsCount    int
}
