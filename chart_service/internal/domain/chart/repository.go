package chart

import (
	"context"

	"github.com/google/uuid"
)

type ChartRepository interface {
	Create(ctx context.Context, c *Chart) (*Chart, error)
	Update(ctx context.Context, c *Chart) (*Chart, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Chart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Chart, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]Chart, error)
}
