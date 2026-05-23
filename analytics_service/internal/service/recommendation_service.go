package service

import (
	"context"
	"math/rand"
	"sort"

	"charts-analytics-service/internal/client"
	"charts-analytics-service/internal/client/dto"
	"charts-analytics-service/internal/domain/reaction"
	"charts-analytics-service/internal/repository/postgres"

	"github.com/google/uuid"
)

type RecommendationService struct {
	repo          *postgres.ReactionRepository
	archiveClient *client.ArchiveClient
}

func NewRecommendationService(
	repo *postgres.ReactionRepository,
	archiveClient *client.ArchiveClient,
) *RecommendationService {
	return &RecommendationService{
		repo:          repo,
		archiveClient: archiveClient,
	}
}

type scoredEpisode struct {
	episode dto.EpisodeResponse
	score   float64
}

func (s *RecommendationService) fallback(
	ctx context.Context,
	limit int,
	ratedCharts map[uuid.UUID]struct{},
) ([]dto.EpisodeResponse, error) {

	fetchLimit := limit * 7

	eps, err := s.archiveClient.GetLatestEpisodes(ctx, fetchLimit)
	if err != nil {
		return nil, err
	}

	if len(eps) == 0 {
		return []dto.EpisodeResponse{}, nil
	}

	var filtered []dto.EpisodeResponse
	for _, ep := range eps {
		if ratedCharts != nil && len(ratedCharts) > 0 {
			if _, ok := ratedCharts[ep.ChartID]; ok {
				continue
			}
		}
		filtered = append(filtered, ep)
	}

	if len(filtered) == 0 {
		return []dto.EpisodeResponse{}, nil
	}

	if len(filtered) <= limit {
		return filtered, nil
	}

	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	return filtered[:limit], nil
}

func (s *RecommendationService) buildUserProfile(
	ctx context.Context,
	reactions []reaction.Reaction,
) (map[string]float64, error) {

	profile := make(map[string]float64)

	for _, r := range reactions {

		ep, err := s.archiveClient.GetNearestLeftEpisode(
			ctx,
			r.ChartID,
			r.CreatedAt,
		)
		if err != nil {
			continue
		}

		if ep == nil {
			continue
		}

		baseWeight := getReactionWeight(r.Type)

		for _, track := range ep.Tracks {

			key := track.Artist + "|" + track.Title
			if key == "" {
				continue
			}

			positionWeight := 1.0 / float64(track.CurrentPosition)

			weight := baseWeight * positionWeight

			profile[key] += weight
		}
	}

	return profile, nil
}

func getReactionWeight(t reaction.ReactionType) float64 {
	switch t {
	case reaction.ReactionLike:
		return 1.0
	case reaction.ReactionView:
		return 0.3
	case reaction.ReactionDislike:
		return -1.0
	default:
		return 0
	}
}

func (s *RecommendationService) scoreEpisode(
	ep dto.EpisodeResponse,
	profile map[string]float64,
) float64 {

	var score float64

	for _, track := range ep.Tracks {

		if track.CurrentPosition <= 0 {
			continue
		}

		key := track.Artist + "|" + track.Title

		weight, ok := profile[key]
		if !ok {
			continue
		}

		positionWeight := 1.0 / float64(track.CurrentPosition)

		score += weight * positionWeight
	}

	return score
}

func keepOnlyTopTrack(
	ep dto.EpisodeResponse,
) dto.EpisodeResponse {

	var topTracks []dto.TrackEpisodeResponse

	for _, track := range ep.Tracks {
		if track.CurrentPosition == 1 {
			topTracks = append(topTracks, track)
			break
		}
	}

	ep.Tracks = topTracks

	return ep
}

func (s *RecommendationService) GetRecommendations(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]dto.EpisodeResponse, error) {

	if userID == uuid.Nil {
		return s.fallback(ctx, limit, nil)
	}

	reactions, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	ratedCharts := make(map[uuid.UUID]struct{})

	for _, r := range reactions {
		ratedCharts[r.ChartID] = struct{}{}
	}

	if len(reactions) == 0 {
		return s.fallback(ctx, limit, ratedCharts)
	}

	profile, err := s.buildUserProfile(ctx, reactions)
	if err != nil {
		return nil, err
	}

	if len(profile) == 0 {
		return s.fallback(ctx, limit, ratedCharts)
	}

	candidates, err := s.archiveClient.GetLatestEpisodes(ctx, limit*10)
	if err != nil {
		return nil, err
	}

	if len(candidates) == 0 {
		return s.fallback(ctx, limit, ratedCharts)
	}

	var scored []scoredEpisode
	const scoreCap = 10.0
	for _, ep := range candidates {

		if _, ok := ratedCharts[ep.ChartID]; ok {
			continue
		}

		score := s.scoreEpisode(ep, profile)

		if score <= 0 {
			continue
		}

		if score > scoreCap {
			score = scoreCap
		}

		scored = append(scored, scoredEpisode{
			episode: ep,
			score:   score,
		})
	}

	if len(scored) == 0 {
		return s.fallback(ctx, limit, ratedCharts)
	}
	bestPerChart := make(map[uuid.UUID]scoredEpisode)

	for _, se := range scored {

		chartID := se.episode.ChartID

		if existing, ok := bestPerChart[chartID]; ok {
			if se.score > existing.score {
				bestPerChart[chartID] = se
			}
		} else {
			bestPerChart[chartID] = se
		}
	}
	var final []scoredEpisode

	for _, v := range bestPerChart {
		final = append(final, v)
	}

	sort.Slice(final, func(i, j int) bool {
		return final[i].score > final[j].score
	})

	if len(final) > limit {
		final = final[:limit]
	}

	result := make([]dto.EpisodeResponse, len(final))

	for i, v := range final {
		result[i] = keepOnlyTopTrack(v.episode)
	}

	return result, nil

}
