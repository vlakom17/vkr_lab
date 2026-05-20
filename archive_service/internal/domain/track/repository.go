package track

import (
	"context"
)

type TrackRepository interface {
	FindOrCreate(ctx context.Context, artist, title, normalizedKey string) (*Track, error)
	Search(ctx context.Context, query string) ([]Track, error)
}
