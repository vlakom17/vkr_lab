package service

import (
	"context"
	"strings"

	"charts-archive-service/internal/domain/track"
	"charts-archive-service/internal/repository/postgres"
)

type TrackService struct {
	trackRepo postgres.TrackRepository
}

func NewTrackService(trackRepo postgres.TrackRepository) *TrackService {
	return &TrackService{
		trackRepo: trackRepo,
	}
}

func (s *TrackService) SearchTracks(ctx context.Context, query string) ([]track.Track, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return []track.Track{}, nil
	}

	return s.trackRepo.Search(ctx, query)
}
