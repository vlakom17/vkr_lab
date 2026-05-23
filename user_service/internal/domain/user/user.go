package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	About        string    `json:"about"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}
