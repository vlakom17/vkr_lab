package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Create(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error
	GetUserID(ctx context.Context, token string) (uuid.UUID, error)
	Delete(ctx context.Context, token string) error
	GetTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error)
}
