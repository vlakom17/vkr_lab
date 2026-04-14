package episode

import (
	"charts-archive-service/internal/repository/dto"
	"context"
	"time"

	"github.com/google/uuid"
)

type EpisodeRepository interface {
	Create(ctx context.Context, episode *Episode) (*Episode, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Episode, error)
	GetByChartID(ctx context.Context, chartID uuid.UUID) ([]Episode, error)
	GetLatestByChartID(ctx context.Context, chartID uuid.UUID) (*Episode, error)
	GetLatestByLimit(ctx context.Context, limit int) ([]Episode, error)
	GetLatestWithTracksByLimit(ctx context.Context, limit int) ([]dto.EpisodeResponse, error)
	GetTracksWithMetaByEpisodeID(ctx context.Context, episodeID uuid.UUID) ([]dto.TrackEpisodeResponse, error)
	GetNearestLeftWithTracks(ctx context.Context, chartID uuid.UUID, date time.Time) (*dto.EpisodeResponse, error)
	GetLatestEpisodesPage(ctx context.Context, limit int, offset int) ([]Episode, error)
}
