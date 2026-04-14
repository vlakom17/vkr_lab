package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"charts-chart-service/internal/client"
	"charts-chart-service/internal/domain/chart"
	"charts-chart-service/internal/domain/event"
	"charts-chart-service/internal/repository/postgres"
	"charts-chart-service/internal/utilits"

	"github.com/google/uuid"
)

type ChartService struct {
	repo            *postgres.ChartRepository
	userClient      *client.UserClient
	analyticsClient *client.AnalyticsClient
	producer        event.EventProducer
}

func NewChartService(
	repo *postgres.ChartRepository,
	userClient *client.UserClient,
	analyticsClient *client.AnalyticsClient,
	producer event.EventProducer,
) *ChartService {
	return &ChartService{
		repo:            repo,
		userClient:      userClient,
		analyticsClient: analyticsClient,
		producer:        producer,
	}
}

func (s *ChartService) CreateChart(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	genre string,
	positionCount int,
	description string,
) (*chart.Chart, error) {

	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(existing) > 0 {
		return nil, errors.New("user already has a chart")
	}

	if title == "" {
		return nil, errors.New("title is required")
	}

	if !chart.IsValidPositionCount(positionCount) {
		return nil, errors.New("invalid position count")
	}

	c := &chart.Chart{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         title,
		Genre:         genre,
		PositionCount: positionCount,
		Description:   description,
	}

	return s.repo.Create(ctx, c)
}

func (s *ChartService) PatchChart(
	ctx context.Context,
	userID uuid.UUID,
	chartID uuid.UUID,
	title *string,
	genre *string,
	positionCount *int,
	description *string,
) (*chart.Chart, error) {

	c, err := s.repo.GetByID(ctx, chartID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("chart not found")
	}

	if c.UserID != userID {
		return nil, errors.New("forbidden")
	}

	if title != nil {
		if *title == "" {
			return nil, errors.New("title cannot be empty")
		}
		c.Title = *title
	}

	if genre != nil {
		c.Genre = *genre
	}

	if positionCount != nil {
		if !chart.IsValidPositionCount(*positionCount) {
			return nil, errors.New("invalid position count")
		}
		c.PositionCount = *positionCount
	}

	if description != nil {
		c.Description = *description
	}

	return s.repo.Update(ctx, c)
}

func (s *ChartService) GetChartByID(
	ctx context.Context,
	id uuid.UUID,
	userID *uuid.UUID,
) (*chart.Chart, error) {

	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errors.New("chart not found")
	}

	if userID != nil && c.UserID != *userID {
		_ = s.producer.SendReaction(ctx, event.ReactionEvent{
			UserID:    *userID,
			ChartID:   c.ID,
			Type:      "view",
			CreatedAt: time.Now(),
		})
	}

	return c, nil
}

func (s *ChartService) GetChartByIDWithoutView(
	ctx context.Context,
	id uuid.UUID,
) (*chart.Chart, error) {

	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errors.New("chart not found")
	}

	return c, nil
}

func (s *ChartService) GetMyChart(
	ctx context.Context,
	userID uuid.UUID,
) (*chart.Chart, error) {

	charts, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(charts) == 0 {
		return nil, errors.New("chart not found")
	}

	return &charts[0], nil
}

func (s *ChartService) GetChartsByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]chart.Chart, error) {
	return s.repo.GetByIDs(ctx, ids)
}

func (s *ChartService) GetChartIDsByGenre(
	ctx context.Context,
	genre string,
	limit int,
) ([]uuid.UUID, error) {
	return s.repo.GetIDsByGenre(ctx, genre, limit)
}

func (s *ChartService) GetGenresByChartIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (map[uuid.UUID]string, error) {
	return s.repo.GetGenresByChartIDs(ctx, ids)
}

func (s *ChartService) SetReaction(
	ctx context.Context,
	userID uuid.UUID,
	chartID uuid.UUID,
	reactionType event.ReactionType,
) error {

	c, err := s.repo.GetByID(ctx, chartID)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("chart not found")
	}

	if reactionType != event.ReactionLike && reactionType != event.ReactionDislike &&
		reactionType != event.ReactionView && reactionType != event.ReactionRemove {
		return errors.New("invalid reaction type")
	}

	ev := event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      reactionType,
		CreatedAt: time.Now(),
	}

	return s.producer.SendReaction(ctx, ev)
}

func (s *ChartService) CreateEpisodeSnapshot(
	ctx context.Context,
	userID uuid.UUID,
	chartID uuid.UUID,
	tracks []event.TrackEntryData,
) error {

	c, err := s.repo.GetByID(ctx, chartID)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("chart not found")
	}

	if c.UserID != userID {
		return errors.New("forbidden")
	}

	if len(tracks) != c.PositionCount {
		return errors.New("invalid number of tracks")
	}

	currentPositions := make(map[int]bool)

	for _, t := range tracks {

		if err := utilits.ValidateTrackInput(t.Artist, t.Title); err != nil {
			return fmt.Errorf("invalid track at position %d: %w", t.CurrentPosition, err)
		}

		if t.CurrentPosition < 1 || t.CurrentPosition > c.PositionCount {
			return errors.New("invalid current position")
		}

		if currentPositions[t.CurrentPosition] {
			return errors.New("duplicate current position")
		}
		currentPositions[t.CurrentPosition] = true

	}

	for i := 1; i <= c.PositionCount; i++ {
		if !currentPositions[i] {
			return errors.New("positions must be continuous")
		}
	}

	ev := event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		Tracks:    tracks,
		CreatedAt: time.Now(),
	}

	return s.producer.SendEpisode(ctx, ev)
}

func (s *ChartService) GetMostPopularCharts(
	ctx context.Context,
	limit int,
) ([]chart.Chart, error) {

	if limit <= 0 {
		limit = 25
	}

	ids, err := s.analyticsClient.GetMostPopularChartIDs(ctx, limit)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []chart.Chart{}, nil
	}

	charts, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	// восстановление порядка
	chartMap := make(map[uuid.UUID]chart.Chart, len(charts))
	for _, c := range charts {
		chartMap[c.ID] = c
	}

	result := make([]chart.Chart, 0, len(ids))
	for _, id := range ids {
		if c, ok := chartMap[id]; ok {
			result = append(result, c)
		}
	}

	return result, nil
}

func (s *ChartService) GetUserLikedCharts(
	ctx context.Context,
	userID uuid.UUID,
) ([]chart.Chart, error) {

	ids, err := s.analyticsClient.GetUserLikedChartIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []chart.Chart{}, nil
	}

	charts, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	chartMap := make(map[uuid.UUID]chart.Chart, len(charts))
	for _, c := range charts {
		chartMap[c.ID] = c
	}

	result := make([]chart.Chart, 0, len(ids))
	for _, id := range ids {
		if c, ok := chartMap[id]; ok {
			result = append(result, c)
		}
	}

	return result, nil
}

func (s *ChartService) GetUserDislikedCharts(
	ctx context.Context,
	userID uuid.UUID,
) ([]chart.Chart, error) {

	ids, err := s.analyticsClient.GetUserDislikedChartIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []chart.Chart{}, nil
	}

	charts, err := s.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	chartMap := make(map[uuid.UUID]chart.Chart, len(charts))
	for _, c := range charts {
		chartMap[c.ID] = c
	}

	result := make([]chart.Chart, 0, len(ids))
	for _, id := range ids {
		if c, ok := chartMap[id]; ok {
			result = append(result, c)
		}
	}

	return result, nil
}
