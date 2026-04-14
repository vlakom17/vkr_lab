package chart

import (
	"time"

	"github.com/google/uuid"
)

type Chart struct {
	ID            uuid.UUID `json:"id,omitempty"`
	UserID        uuid.UUID `json:"user_id"`
	Title         string    `json:"title"`
	Genre         string    `json:"genre"`
	Description   string    `json:"description"`
	PositionCount int       `json:"position_count"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}
