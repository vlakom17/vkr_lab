package service

import (
	"context"
	"fmt"
	"time"

	"charts-archive-service/internal/domain/episode"
	"charts-archive-service/internal/domain/event"
	"charts-archive-service/internal/domain/track_episode"
	"charts-archive-service/internal/repository/dto"
	"charts-archive-service/internal/repository/postgres"
	"charts-archive-service/internal/transport/http/response"
	"charts-archive-service/internal/utilits"

	"github.com/google/uuid"
)

type EpisodeService struct {
	episodeRepo postgres.EpisodeRepository
	trackRepo   postgres.TrackRepository
}

func NewEpisodeService(
	episodeRepo postgres.EpisodeRepository,
	trackRepo postgres.TrackRepository,
) *EpisodeService {
	return &EpisodeService{
		episodeRepo: episodeRepo,
		trackRepo:   trackRepo,
	}
}

func (s *EpisodeService) GetEpisode(
	ctx context.Context,
	id uuid.UUID,
) (*response.EpisodeResponse, error) {

	ep, err := s.episodeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ep == nil {
		return nil, nil
	}

	tracks, err := s.episodeRepo.GetTracksWithMetaByEpisodeID(ctx, id)
	if err != nil {
		return nil, err
	}

	respTracks := make([]response.TrackEpisodeResponse, 0, len(tracks))

	for _, t := range tracks {
		links := utilits.BuildListenLinks(t.Artist, t.Title)

		respTracks = append(respTracks, response.TrackEpisodeResponse{
			TrackID:             t.TrackID,
			Artist:              utilits.Capitalize(t.Artist),
			Title:               utilits.Capitalize(t.Title),
			CurrentPosition:     t.CurrentPosition,
			PreviousPosition:    t.PreviousPosition,
			HighestPosition:     t.HighestPosition,
			TimesAtPeakPosition: t.TimesAtPeakPosition,
			EpisodesCount:       t.EpisodesCount,
			ListenLinks: response.ListenLinks{
				AppleMusic:  links.AppleMusic,
				YandexMusic: links.YandexMusic,
			},
		})
	}

	resp := &response.EpisodeResponse{
		ID:        ep.ID,
		ChartID:   ep.ChartID,
		CreatedAt: ep.CreatedAt,
		Tracks:    respTracks,
	}

	return resp, nil
}

func (s *EpisodeService) GetEpisodesByChart(ctx context.Context, chartID uuid.UUID) ([]episode.Episode, error) {
	return s.episodeRepo.GetByChartID(ctx, chartID)
}

func (s *EpisodeService) GetLatestEpisodesWithTracks(
	ctx context.Context,
	limit int,
) ([]dto.EpisodeResponse, error) {
	return s.episodeRepo.GetLatestWithTracksByLimit(ctx, limit)
}

func (s *EpisodeService) HandleEpisodeCreatedEvent(
	ctx context.Context,
	e event.EpisodeSnapshotEvent,
) error {

	prevEpisode, err := s.episodeRepo.GetLatestByChartID(ctx, e.ChartID)
	if err != nil {
		return err
	}

	prevMap := make(map[uuid.UUID]track_episode.TrackEpisode)

	if prevEpisode != nil {
		for _, t := range prevEpisode.Tracks {
			prevMap[t.TrackID] = t
		}
	}

	newEpisode := &episode.Episode{
		ID:        uuid.New(),
		ChartID:   e.ChartID,
		CreatedAt: e.CreatedAt,
	}

	seenTracks := make(map[string]struct{})
	for _, incoming := range e.Tracks {

		nt := utilits.NormalizeTrack(incoming.Artist, incoming.Title)

		if _, exists := seenTracks[nt.NormalizedKey]; exists {
			return fmt.Errorf("duplicate track in episode: %s", nt.NormalizedKey)
		}

		seenTracks[nt.NormalizedKey] = struct{}{}

		track, err := s.trackRepo.FindOrCreate(
			ctx,
			nt.Artist,
			nt.Title,
			nt.NormalizedKey,
		)

		if err != nil {
			return err
		}

		prev, exists := prevMap[track.ID]

		var previousPosition int
		var highestPosition int
		var timesAtPeak int
		var episodesCount int

		if exists {
			previousPosition = prev.CurrentPosition
			episodesCount = prev.EpisodesCount + 1

			if incoming.CurrentPosition < prev.HighestPosition {
				highestPosition = incoming.CurrentPosition
				timesAtPeak = 1

			} else if incoming.CurrentPosition == prev.HighestPosition {
				highestPosition = prev.HighestPosition
				timesAtPeak = prev.TimesAtPeakPosition + 1

			} else {
				highestPosition = prev.HighestPosition
				timesAtPeak = prev.TimesAtPeakPosition
			}

		} else {
			previousPosition = 0
			highestPosition = incoming.CurrentPosition
			episodesCount = 1
			timesAtPeak = 1
		}

		newEpisode.Tracks = append(newEpisode.Tracks, track_episode.TrackEpisode{
			ID:                  uuid.New(),
			EpisodeID:           newEpisode.ID,
			TrackID:             track.ID,
			CurrentPosition:     incoming.CurrentPosition,
			PreviousPosition:    previousPosition,
			HighestPosition:     highestPosition,
			TimesAtPeakPosition: timesAtPeak,
			EpisodesCount:       episodesCount,
		})
	}
	_, err = s.episodeRepo.Create(ctx, newEpisode)

	if err != nil {
		return err
	}

	return nil
}

func (s *EpisodeService) GetNearestLeftEpisode(
	ctx context.Context,
	chartID uuid.UUID,
	date time.Time,
) (*dto.EpisodeResponse, error) {
	return s.episodeRepo.GetNearestLeftWithTracks(ctx, chartID, date)
}

func (s *EpisodeService) GetLatestEpisodesPage(
	ctx context.Context,
	page int,
	limit int,
) ([]episode.Episode, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	return s.episodeRepo.GetLatestEpisodesPage(ctx, limit, offset)
}
