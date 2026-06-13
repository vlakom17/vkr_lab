package service

import (
	"context"
	"strings"

	"charts-archive-service/internal/domain/track"
)

type TrackService struct {
	trackRepo track.TrackRepository
}

func NewTrackService(trackRepo track.TrackRepository) *TrackService {
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
