package auth

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
