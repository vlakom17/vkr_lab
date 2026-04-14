package track

import (
	"context"

	"github.com/google/uuid"
)

type TrackRepository interface {
	FindOrCreate(ctx context.Context, artist, title, normalizedKey string) (*Track, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Track, error)
	Search(ctx context.Context, query string) ([]Track, error)
}
