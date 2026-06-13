package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"charts-analytics-service/internal/client/dto"
	"charts-analytics-service/internal/domain/reaction"

	"github.com/google/uuid"
)

type fakeRecommendationRepository struct {
	reactions      []reaction.Reaction
	getByUserIDErr error
}

func (r *fakeRecommendationRepository) Upsert(ctx context.Context, reaction *reaction.Reaction) (*reaction.Reaction, error) {
	return nil, nil
}

func (r *fakeRecommendationRepository) GetByUserAndChart(ctx context.Context, userID, chartID uuid.UUID) (*reaction.Reaction, error) {
	return nil, nil
}

func (r *fakeRecommendationRepository) CountByChart(ctx context.Context, chartID uuid.UUID) (likes, dislikes, views int, err error) {
	return 0, 0, 0, nil
}

func (r *fakeRecommendationRepository) GetMostPopularChartIDs(ctx context.Context, limit int) ([]uuid.UUID, error) {
	return nil, nil
}

func (r *fakeRecommendationRepository) GetUserChartIDsByType(ctx context.Context, userID uuid.UUID, t reaction.ReactionType) ([]uuid.UUID, error) {
	return nil, nil
}

func (r *fakeRecommendationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]reaction.Reaction, error) {
	if r.getByUserIDErr != nil {
		return nil, r.getByUserIDErr
	}

	return r.reactions, nil
}

type fakeArchiveClient struct {
	latestEpisodes []dto.EpisodeResponse
	latestLimit    int
	latestErr      error

	nearestEpisodes map[uuid.UUID]*dto.EpisodeResponse
	nearestErr      error
}

func (c *fakeArchiveClient) GetNearestLeftEpisode(
	ctx context.Context,
	chartID uuid.UUID,
	date time.Time,
) (*dto.EpisodeResponse, error) {
	if c.nearestErr != nil {
		return nil, c.nearestErr
	}

	if c.nearestEpisodes == nil {
		return nil, nil
	}

	return c.nearestEpisodes[chartID], nil
}

func (c *fakeArchiveClient) GetLatestEpisodes(
	ctx context.Context,
	limit int,
) ([]dto.EpisodeResponse, error) {
	c.latestLimit = limit

	if c.latestErr != nil {
		return nil, c.latestErr
	}

	return c.latestEpisodes, nil
}

func recommendationEpisode(chartID uuid.UUID, artist string, title string, position int) dto.EpisodeResponse {
	return dto.EpisodeResponse{
		ID:      uuid.New(),
		ChartID: chartID,
		Tracks: []dto.TrackEpisodeResponse{
			{
				TrackID:         uuid.New(),
				Artist:          artist,
				Title:           title,
				CurrentPosition: position,
			},
		},
	}
}
func TestGetReactionWeight(t *testing.T) {
	tests := []struct {
		name     string
		input    reaction.ReactionType
		expected float64
	}{
		{"like", reaction.ReactionLike, 1.0},
		{"view", reaction.ReactionView, 0.3},
		{"dislike", reaction.ReactionDislike, -1.0},
		{"unknown", reaction.ReactionType("unknown"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getReactionWeight(tt.input)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestScoreEpisode_UsesProfileAndPositionWeight(t *testing.T) {
	service := NewRecommendationService(
		&fakeRecommendationRepository{},
		&fakeArchiveClient{},
	)

	chartID := uuid.New()

	ep := dto.EpisodeResponse{
		ID:      uuid.New(),
		ChartID: chartID,
		Tracks: []dto.TrackEpisodeResponse{
			{
				Artist:          "artist",
				Title:           "song",
				CurrentPosition: 1,
			},
			{
				Artist:          "artist 2",
				Title:           "song 2",
				CurrentPosition: 2,
			},
		},
	}

	profile := map[string]float64{
		"artist|song":     1.0,
		"artist 2|song 2": 0.5,
	}

	result := service.scoreEpisode(ep, profile)

	expected := 1.25

	if result != expected {
		t.Errorf("expected score %v, got %v", expected, result)
	}
}

func TestKeepOnlyTopTrack_KeepsOnlyPositionOne(t *testing.T) {
	chartID := uuid.New()

	ep := dto.EpisodeResponse{
		ID:      uuid.New(),
		ChartID: chartID,
		Tracks: []dto.TrackEpisodeResponse{
			{Artist: "artist 2", Title: "song 2", CurrentPosition: 2},
			{Artist: "artist 1", Title: "song 1", CurrentPosition: 1},
			{Artist: "artist 3", Title: "song 3", CurrentPosition: 3},
		},
	}

	result := keepOnlyTopTrack(ep)

	if len(result.Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(result.Tracks))
	}

	if result.Tracks[0].CurrentPosition != 1 {
		t.Errorf("expected position 1, got %d", result.Tracks[0].CurrentPosition)
	}

	if result.Tracks[0].Artist != "artist 1" {
		t.Errorf("expected artist %q, got %q", "artist 1", result.Tracks[0].Artist)
	}
}

func TestFallback_ReturnsLatestEpisodesWhenUserIsAnonymous(t *testing.T) {
	chartID := uuid.New()

	archiveClient := &fakeArchiveClient{
		latestEpisodes: []dto.EpisodeResponse{
			recommendationEpisode(chartID, "artist", "song", 1),
		},
	}

	service := NewRecommendationService(
		&fakeRecommendationRepository{},
		archiveClient,
	)

	result, err := service.GetRecommendations(context.Background(), uuid.Nil, 3)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if archiveClient.latestLimit != 21 {
		t.Errorf("expected latest limit 21, got %d", archiveClient.latestLimit)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(result))
	}
}

func TestGetRecommendations_ReturnsScoredCandidateFromLikedTrack(t *testing.T) {
	userID := uuid.New()

	likedChartID := uuid.New()
	recommendedChartID := uuid.New()

	now := time.Now()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   likedChartID,
				Type:      reaction.ReactionLike,
				CreatedAt: now,
			},
		},
	}

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{
			likedChartID: {
				ID:      uuid.New(),
				ChartID: likedChartID,
				Tracks: []dto.TrackEpisodeResponse{
					{
						Artist:          "artist",
						Title:           "song",
						CurrentPosition: 1,
					},
				},
			},
		},
		latestEpisodes: []dto.EpisodeResponse{
			{
				ID:      uuid.New(),
				ChartID: recommendedChartID,
				Tracks: []dto.TrackEpisodeResponse{
					{
						Artist:          "artist",
						Title:           "song",
						CurrentPosition: 1,
					},
					{
						Artist:          "other",
						Title:           "track",
						CurrentPosition: 2,
					},
				},
			},
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(context.Background(), userID, 5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(result))
	}

	if result[0].ChartID != recommendedChartID {
		t.Errorf("expected chartID %s, got %s", recommendedChartID, result[0].ChartID)
	}

	if len(result[0].Tracks) != 1 {
		t.Fatalf("expected only top track, got %d tracks", len(result[0].Tracks))
	}

	if result[0].Tracks[0].CurrentPosition != 1 {
		t.Errorf("expected top track position 1, got %d", result[0].Tracks[0].CurrentPosition)
	}
}
func TestGetRecommendations_ExcludesRatedCharts(t *testing.T) {
	userID := uuid.New()

	likedChartID := uuid.New()
	candidateChartID := uuid.New()

	now := time.Now()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   likedChartID,
				Type:      reaction.ReactionLike,
				CreatedAt: now,
			},
		},
	}

	likedEpisode := recommendationEpisode(
		likedChartID,
		"artist",
		"song",
		1,
	)

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{
			likedChartID: &likedEpisode,
		},
		latestEpisodes: []dto.EpisodeResponse{
			recommendationEpisode(likedChartID, "artist", "song", 1),
			recommendationEpisode(candidateChartID, "artist", "song", 1),
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(context.Background(), userID, 5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(result))
	}

	if result[0].ChartID != candidateChartID {
		t.Errorf("expected only unrated chart %s, got %s", candidateChartID, result[0].ChartID)
	}
}
func TestGetRecommendations_FallsBackWhenOnlyNegativeScores(t *testing.T) {
	userID := uuid.New()

	dislikedChartID := uuid.New()
	candidateChartID := uuid.New()
	fallbackChartID := uuid.New()

	now := time.Now()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   dislikedChartID,
				Type:      reaction.ReactionDislike,
				CreatedAt: now,
			},
		},
	}

	dislikedEpisode := recommendationEpisode(
		dislikedChartID,
		"artist",
		"song",
		1,
	)

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{
			dislikedChartID: &dislikedEpisode,
		},
		latestEpisodes: []dto.EpisodeResponse{
			recommendationEpisode(candidateChartID, "artist", "song", 1),
			recommendationEpisode(fallbackChartID, "other", "track", 1),
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(context.Background(), userID, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 fallback recommendation, got %d", len(result))
	}
}

func TestGetRecommendations_KeepBestEpisodePerChart(t *testing.T) {
	userID := uuid.New()

	likedChartID := uuid.New()
	targetChartID := uuid.New()

	now := time.Now()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   likedChartID,
				Type:      reaction.ReactionLike,
				CreatedAt: now,
			},
		},
	}

	likedEpisode := recommendationEpisode(
		likedChartID,
		"artist",
		"song",
		1,
	)

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{
			likedChartID: &likedEpisode,
		},
		latestEpisodes: []dto.EpisodeResponse{
			{
				ID:      uuid.New(),
				ChartID: targetChartID,
				Tracks: []dto.TrackEpisodeResponse{
					{
						Artist:          "artist",
						Title:           "song",
						CurrentPosition: 1,
					},
				},
			},
			{
				ID:      uuid.New(),
				ChartID: targetChartID,
				Tracks: []dto.TrackEpisodeResponse{
					{
						Artist:          "artist",
						Title:           "song",
						CurrentPosition: 2,
					},
				},
			},
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(
		context.Background(),
		userID,
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 episode, got %d", len(result))
	}

	if result[0].ChartID != targetChartID {
		t.Errorf("unexpected chart id")
	}
}
func TestGetRecommendations_SortsByScoreDescending(t *testing.T) {
	userID := uuid.New()

	likedChartID := uuid.New()
	chartA := uuid.New()
	chartB := uuid.New()

	now := time.Now()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   likedChartID,
				Type:      reaction.ReactionLike,
				CreatedAt: now,
			},
		},
	}

	likedEpisode := recommendationEpisode(
		likedChartID,
		"artist",
		"song",
		1,
	)

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{
			likedChartID: &likedEpisode,
		},
		latestEpisodes: []dto.EpisodeResponse{
			{
				ID:      uuid.New(),
				ChartID: chartA,
				Tracks: []dto.TrackEpisodeResponse{
					{
						Artist:          "artist",
						Title:           "song",
						CurrentPosition: 2,
					},
				},
			},
			{
				ID:      uuid.New(),
				ChartID: chartB,
				Tracks: []dto.TrackEpisodeResponse{
					{
						Artist:          "artist",
						Title:           "song",
						CurrentPosition: 1,
					},
				},
			},
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(
		context.Background(),
		userID,
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(result))
	}

	if result[0].ChartID != chartB {
		t.Errorf("expected highest scored chart first")
	}
}
func TestGetRecommendations_ReturnsErrorWhenGetByUserIDFails(t *testing.T) {
	userID := uuid.New()
	expectedErr := errors.New("database error")

	repo := &fakeRecommendationRepository{
		getByUserIDErr: expectedErr,
	}

	service := NewRecommendationService(repo, &fakeArchiveClient{})

	result, err := service.GetRecommendations(context.Background(), userID, 5)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
func TestGetRecommendations_UsesFallbackWhenUserHasNoReactions(t *testing.T) {
	userID := uuid.New()
	fallbackChartID := uuid.New()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{},
	}

	archiveClient := &fakeArchiveClient{
		latestEpisodes: []dto.EpisodeResponse{
			recommendationEpisode(fallbackChartID, "artist", "song", 1),
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(context.Background(), userID, 3)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if archiveClient.latestLimit != 21 {
		t.Errorf("expected fallback fetch limit 21, got %d", archiveClient.latestLimit)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 fallback recommendation, got %d", len(result))
	}

	if result[0].ChartID != fallbackChartID {
		t.Errorf("expected fallback chart %s, got %s", fallbackChartID, result[0].ChartID)
	}
}
func TestGetRecommendations_UsesFallbackWhenProfileIsEmpty(t *testing.T) {
	userID := uuid.New()
	ratedChartID := uuid.New()
	fallbackChartID := uuid.New()

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   ratedChartID,
				Type:      reaction.ReactionLike,
				CreatedAt: time.Now(),
			},
		},
	}

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{},
		latestEpisodes: []dto.EpisodeResponse{
			recommendationEpisode(fallbackChartID, "artist", "song", 1),
		},
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(context.Background(), userID, 3)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 fallback recommendation, got %d", len(result))
	}

	if result[0].ChartID != fallbackChartID {
		t.Errorf("expected fallback chart %s, got %s", fallbackChartID, result[0].ChartID)
	}
}
func TestGetRecommendations_ReturnsErrorWhenGetLatestEpisodesFails(t *testing.T) {
	userID := uuid.New()
	likedChartID := uuid.New()

	expectedErr := errors.New("archive error")

	likedEpisode := recommendationEpisode(likedChartID, "artist", "song", 1)

	repo := &fakeRecommendationRepository{
		reactions: []reaction.Reaction{
			{
				UserID:    userID,
				ChartID:   likedChartID,
				Type:      reaction.ReactionLike,
				CreatedAt: time.Now(),
			},
		},
	}

	archiveClient := &fakeArchiveClient{
		nearestEpisodes: map[uuid.UUID]*dto.EpisodeResponse{
			likedChartID: &likedEpisode,
		},
		latestErr: expectedErr,
	}

	service := NewRecommendationService(repo, archiveClient)

	result, err := service.GetRecommendations(context.Background(), userID, 3)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
func TestFallback_ExcludesRatedCharts(t *testing.T) {
	ratedChartID := uuid.New()
	unratedChartID := uuid.New()

	archiveClient := &fakeArchiveClient{
		latestEpisodes: []dto.EpisodeResponse{
			recommendationEpisode(ratedChartID, "artist", "song", 1),
			recommendationEpisode(unratedChartID, "other", "track", 1),
		},
	}

	service := NewRecommendationService(
		&fakeRecommendationRepository{},
		archiveClient,
	)

	ratedCharts := map[uuid.UUID]struct{}{
		ratedChartID: {},
	}

	result, err := service.fallback(context.Background(), 5, ratedCharts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 fallback episode, got %d", len(result))
	}

	if result[0].ChartID != unratedChartID {
		t.Errorf("expected only unrated chart %s, got %s", unratedChartID, result[0].ChartID)
	}
}
