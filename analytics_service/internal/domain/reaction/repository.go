package reaction

import (
	"context"

	"github.com/google/uuid"
)

type ReactionRepository interface {
	Upsert(ctx context.Context, reaction *Reaction) (*Reaction, error)
	GetByUserAndChart(ctx context.Context, userID, chartID uuid.UUID) (*Reaction, error)
	CountByChart(ctx context.Context, chartID uuid.UUID) (likes, dislikes, views int, err error)
	GetMostPopularChartIDs(ctx context.Context, limit int) ([]uuid.UUID, error)
	GetUserChartIDsByType(ctx context.Context, userID uuid.UUID, t ReactionType) ([]uuid.UUID, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Reaction, error)
}
