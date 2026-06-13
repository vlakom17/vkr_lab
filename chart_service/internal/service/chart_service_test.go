package service

import (
	"context"
	"errors"
	"testing"

	"charts-chart-service/internal/domain/chart"
	"charts-chart-service/internal/domain/event"

	"github.com/google/uuid"
)

type fakeChartRepository struct {
	charts         []chart.Chart
	chartByID      *chart.Chart
	createdChart   *chart.Chart
	updatedChart   *chart.Chart
	getByUserIDErr error
	createErr      error
	getByIDErr     error
	updateErr      error
}

func (r *fakeChartRepository) Create(ctx context.Context, c *chart.Chart) (*chart.Chart, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}

	r.createdChart = c
	r.charts = append(r.charts, *c)

	return c, nil
}

func (r *fakeChartRepository) Update(ctx context.Context, c *chart.Chart) (*chart.Chart, error) {
	if r.updateErr != nil {
		return nil, r.updateErr
	}

	r.updatedChart = c
	r.chartByID = c

	return c, nil
}

func (r *fakeChartRepository) GetByID(ctx context.Context, id uuid.UUID) (*chart.Chart, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}

	return r.chartByID, nil
}

func (r *fakeChartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]chart.Chart, error) {
	if r.getByUserIDErr != nil {
		return nil, r.getByUserIDErr
	}

	return r.charts, nil
}

func (r *fakeChartRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]chart.Chart, error) {
	return r.charts, nil
}

type fakeAnalyticsClient struct {
	popularIDs  []uuid.UUID
	likedIDs    []uuid.UUID
	dislikedIDs []uuid.UUID

	limit int
	err   error
}

func (c *fakeAnalyticsClient) GetMostPopularChartIDs(ctx context.Context, limit int) ([]uuid.UUID, error) {
	c.limit = limit

	if c.err != nil {
		return nil, c.err
	}

	return c.popularIDs, nil
}

func (c *fakeAnalyticsClient) GetUserLikedChartIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.likedIDs, nil
}

func (c *fakeAnalyticsClient) GetUserDislikedChartIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.dislikedIDs, nil
}

type fakeEventProducer struct {
	reactionEvent      *event.ReactionEvent
	episodeEvent       *event.EpisodeSnapshotEvent
	sendReactionCalled bool
	sendEpisodeCalled  bool
	sendEpisodeErr     error
	sendReactionErr    error
}

func (p *fakeEventProducer) SendReaction(ctx context.Context, ev event.ReactionEvent) error {
	p.sendReactionCalled = true
	p.reactionEvent = &ev

	if p.sendReactionErr != nil {
		return p.sendReactionErr
	}

	return nil
}

func (p *fakeEventProducer) SendEpisode(ctx context.Context, ev event.EpisodeSnapshotEvent) error {
	p.sendEpisodeCalled = true
	p.episodeEvent = &ev

	if p.sendEpisodeErr != nil {
		return p.sendEpisodeErr
	}

	return nil
}

func TestCreateChart_CreatesChartWhenUserHasNoChart(t *testing.T) {
	userID := uuid.New()

	repo := &fakeChartRepository{}
	analyticsClient := &fakeAnalyticsClient{}
	producer := &fakeEventProducer{}

	service := NewChartService(repo, nil, analyticsClient, producer)

	result, err := service.CreateChart(
		context.Background(),
		userID,
		"My Chart",
		"pop",
		"description",
		10,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected chart, got nil")
	}

	if result.ID == uuid.Nil {
		t.Errorf("expected generated chart ID")
	}

	if result.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, result.UserID)
	}

	if result.Title != "My Chart" {
		t.Errorf("expected title %s, got %s", "My Chart", result.Title)
	}

	if result.Genre != "pop" {
		t.Errorf("expected genre %s, got %s", "pop", result.Genre)
	}

	if result.Description != "description" {
		t.Errorf("expected description %s, got %s", "description", result.Description)
	}

	if result.PositionCount != 10 {
		t.Errorf("expected position count %d, got %d", 10, result.PositionCount)
	}

	if repo.createdChart != result {
		t.Errorf("expected chart to be saved in repository")
	}
}

func TestCreateChart_ReturnsErrorWhenUserAlreadyHasChart(t *testing.T) {
	userID := uuid.New()

	repo := &fakeChartRepository{
		charts: []chart.Chart{
			{
				ID:            uuid.New(),
				UserID:        userID,
				Title:         "Existing Chart",
				Genre:         "Pop",
				Description:   "Already exists",
				PositionCount: 10,
			},
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	result, err := service.CreateChart(
		context.Background(),
		userID,
		"New Chart",
		"Rock",
		"Description",
		10,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.createdChart != nil {
		t.Errorf("Create should not be called")
	}
}

func TestCreateChart_ReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	userID := uuid.New()

	repo := &fakeChartRepository{}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	result, err := service.CreateChart(
		context.Background(),
		userID,
		"",
		"Pop",
		"Description",
		10,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.createdChart != nil {
		t.Errorf("Create should not be called when title is empty")
	}
}

func TestCreateChart_ReturnsErrorWhenPositionCountIsInvalid(t *testing.T) {
	userID := uuid.New()

	repo := &fakeChartRepository{}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	result, err := service.CreateChart(
		context.Background(),
		userID,
		"My Chart",
		"Pop",
		"Description",
		123,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.createdChart != nil {
		t.Errorf("Create should not be called when position count is invalid")
	}
}

func TestCreateChart_ReturnsErrorWhenGetByUserIDFails(t *testing.T) {
	userID := uuid.New()

	expectedErr := errors.New("database error")

	repo := &fakeChartRepository{
		getByUserIDErr: expectedErr,
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	result, err := service.CreateChart(
		context.Background(),
		userID,
		"My Chart",
		"Pop",
		"Description",
		10,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.createdChart != nil {
		t.Errorf("Create should not be called when GetByUserID fails")
	}
}
func TestCreateChart_ReturnsErrorWhenCreateFails(t *testing.T) {
	userID := uuid.New()

	expectedErr := errors.New("create failed")

	repo := &fakeChartRepository{
		createErr: expectedErr,
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	result, err := service.CreateChart(
		context.Background(),
		userID,
		"My Chart",
		"Pop",
		"Description",
		10,
	)

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

func TestPatchChart_UpdatesChartWhenUserIsOwner(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "Old Title",
			Genre:         "Old Genre",
			Description:   "Old Description",
			PositionCount: 10,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	newTitle := "New Title"
	newGenre := "New Genre"
	newDescription := "New Description"
	newPositionCount := 20

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		userID,
		&newTitle,
		&newGenre,
		&newDescription,
		&newPositionCount,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Title != newTitle {
		t.Errorf("expected title %s, got %s", newTitle, result.Title)
	}

	if result.Genre != newGenre {
		t.Errorf("expected genre %s, got %s", newGenre, result.Genre)
	}

	if result.Description != newDescription {
		t.Errorf("expected description %s, got %s", newDescription, result.Description)
	}

	if result.PositionCount != newPositionCount {
		t.Errorf("expected position count %d, got %d", newPositionCount, result.PositionCount)
	}

	if repo.updatedChart != result {
		t.Errorf("expected Update to be called with result chart")
	}
}

func TestPatchChart_ReturnsErrorWhenChartNotFound(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeChartRepository{
		chartByID: nil,
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	newTitle := "New Title"

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		userID,
		&newTitle,
		nil,
		nil,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.updatedChart != nil {
		t.Errorf("Update should not be called when chart not found")
	}
}

func TestPatchChart_ReturnsErrorWhenUserIsNotOwner(t *testing.T) {
	ownerID := uuid.New()
	anotherUserID := uuid.New()
	chartID := uuid.New()

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        ownerID,
			Title:         "Old Title",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 10,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	newTitle := "Hacked Title"

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		anotherUserID,
		&newTitle,
		nil,
		nil,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.updatedChart != nil {
		t.Errorf("Update should not be called when user is not owner")
	}
}

func TestPatchChart_ReturnsErrorWhenTitleIsEmpty(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "Old Title",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 10,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	emptyTitle := ""

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		userID,
		&emptyTitle,
		nil,
		nil,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.updatedChart != nil {
		t.Errorf("Update should not be called when title is empty")
	}
}
func TestPatchChart_ReturnsErrorWhenPositionCountIsInvalid(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "Old Title",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 10,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	invalidPositionCount := 123

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		userID,
		nil,
		nil,
		nil,
		&invalidPositionCount,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.updatedChart != nil {
		t.Errorf("Update should not be called when position count is invalid")
	}
}

func TestPatchChart_ReturnsErrorWhenGetByIDFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("database error")

	repo := &fakeChartRepository{
		getByIDErr: expectedErr,
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	newTitle := "New Title"

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		userID,
		&newTitle,
		nil,
		nil,
		nil,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if repo.updatedChart != nil {
		t.Errorf("Update should not be called when GetByID fails")
	}
}
func TestPatchChart_ReturnsErrorWhenUpdateFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("update failed")

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "Old Title",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 10,
		},
		updateErr: expectedErr,
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		&fakeEventProducer{},
	)

	newTitle := "New Title"

	result, err := service.PatchChart(
		context.Background(),
		chartID,
		userID,
		&newTitle,
		nil,
		nil,
		nil,
	)

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

func TestCreateEpisodeSnapshot_SendsEpisodeWhenInputIsValid(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	producer := &fakeEventProducer{}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		producer,
	)

	tracks := []event.TrackEntryData{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.CreateEpisodeSnapshot(
		context.Background(),
		userID,
		chartID,
		tracks,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !producer.sendEpisodeCalled {
		t.Fatalf("expected SendEpisode to be called")
	}

	if producer.episodeEvent == nil {
		t.Fatalf("expected episode event, got nil")
	}

	if producer.episodeEvent.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, producer.episodeEvent.ChartID)
	}

	if len(producer.episodeEvent.Tracks) != len(tracks) {
		t.Errorf("expected %d tracks, got %d", len(tracks), len(producer.episodeEvent.Tracks))
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenUserIsNotOwner(t *testing.T) {
	ownerID := uuid.New()
	anotherUserID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        ownerID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		producer,
	)

	tracks := []event.TrackEntryData{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.CreateEpisodeSnapshot(
		context.Background(),
		anotherUserID,
		chartID,
		tracks,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendEpisodeCalled {
		t.Errorf("SendEpisode should not be called when user is not owner")
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenTrackCountIsInvalid(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		producer,
	)

	tracks := []event.TrackEntryData{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
	}

	err := service.CreateEpisodeSnapshot(
		context.Background(),
		userID,
		chartID,
		tracks,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendEpisodeCalled {
		t.Errorf("SendEpisode should not be called when track count is invalid")
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenPositionIsOutOfRange(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		producer,
	)

	tracks := []event.TrackEntryData{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 6},
	}

	err := service.CreateEpisodeSnapshot(
		context.Background(),
		userID,
		chartID,
		tracks,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendEpisodeCalled {
		t.Errorf("SendEpisode should not be called when position is out of range")
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenPositionIsDuplicated(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		producer,
	)

	tracks := []event.TrackEntryData{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 2},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 3},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 4},
	}

	err := service.CreateEpisodeSnapshot(
		context.Background(),
		userID,
		chartID,
		tracks,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendEpisodeCalled {
		t.Errorf("SendEpisode should not be called when position is duplicated")
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenTrackInputIsInvalid(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	tracks := []event.TrackEntryData{
		{Artist: "", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.CreateEpisodeSnapshot(context.Background(), userID, chartID, tracks)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendEpisodeCalled {
		t.Errorf("SendEpisode should not be called when track input is invalid")
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenChartNotFound(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: nil,
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	err := service.CreateEpisodeSnapshot(
		context.Background(),
		userID,
		chartID,
		[]event.TrackEntryData{},
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendEpisodeCalled {
		t.Errorf("SendEpisode should not be called when chart not found")
	}
}

func TestCreateEpisodeSnapshot_ReturnsErrorWhenSendEpisodeFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("kafka error")

	producer := &fakeEventProducer{
		sendEpisodeErr: expectedErr,
	}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        userID,
			Title:         "My Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	tracks := []event.TrackEntryData{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.CreateEpisodeSnapshot(context.Background(), userID, chartID, tracks)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
func TestSetReaction_SendsReactionWhenInputIsValid(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        uuid.New(),
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	err := service.SetReaction(
		context.Background(),
		userID,
		chartID,
		event.ReactionLike,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !producer.sendReactionCalled {
		t.Fatalf("expected SendReaction to be called")
	}

	if producer.reactionEvent == nil {
		t.Fatalf("expected reaction event, got nil")
	}

	if producer.reactionEvent.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, producer.reactionEvent.UserID)
	}

	if producer.reactionEvent.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, producer.reactionEvent.ChartID)
	}

	if producer.reactionEvent.Type != event.ReactionLike {
		t.Errorf("expected reaction type %s, got %s", event.ReactionLike, producer.reactionEvent.Type)
	}
}
func TestSetReaction_ReturnsErrorWhenReactionTypeIsInvalid(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        uuid.New(),
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	err := service.SetReaction(
		context.Background(),
		userID,
		chartID,
		event.ReactionType("bad-type"),
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called when reaction type is invalid")
	}
}

func TestSetReaction_ReturnsErrorWhenChartNotFound(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: nil,
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	err := service.SetReaction(
		context.Background(),
		userID,
		chartID,
		event.ReactionLike,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called when chart not found")
	}
}

func TestSetReaction_ReturnsErrorWhenSendReactionFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("kafka error")

	producer := &fakeEventProducer{
		sendReactionErr: expectedErr,
	}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        uuid.New(),
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(
		repo,
		nil,
		&fakeAnalyticsClient{},
		producer,
	)

	err := service.SetReaction(
		context.Background(),
		userID,
		chartID,
		event.ReactionLike,
	)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
func TestGetChartByID_SendsViewWhenViewerIsNotOwner(t *testing.T) {
	ownerID := uuid.New()
	viewerID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        ownerID,
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByID(context.Background(), chartID, &viewerID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected chart, got nil")
	}

	if !producer.sendReactionCalled {
		t.Fatalf("expected SendReaction to be called")
	}

	if producer.reactionEvent.Type != event.ReactionView {
		t.Errorf("expected reaction type view, got %s", producer.reactionEvent.Type)
	}

	if producer.reactionEvent.UserID != viewerID {
		t.Errorf("expected userID %s, got %s", viewerID, producer.reactionEvent.UserID)
	}

	if producer.reactionEvent.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, producer.reactionEvent.ChartID)
	}
}
func TestGetChartByID_DoesNotSendViewWhenViewerIsOwner(t *testing.T) {
	ownerID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        ownerID,
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByID(context.Background(), chartID, &ownerID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected chart, got nil")
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called when owner views own chart")
	}
}
func TestGetChartByID_DoesNotSendViewWhenUserIDIsNil(t *testing.T) {
	ownerID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        ownerID,
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByID(context.Background(), chartID, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected chart, got nil")
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called when userID is nil")
	}
}
func TestGetChartByID_ReturnsErrorWhenChartNotFound(t *testing.T) {
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: nil,
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByID(context.Background(), chartID, nil)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called when chart not found")
	}
}
func TestGetChartByID_IgnoresSendViewError(t *testing.T) {
	ownerID := uuid.New()
	viewerID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{
		sendReactionErr: errors.New("kafka error"),
	}

	repo := &fakeChartRepository{
		chartByID: &chart.Chart{
			ID:            chartID,
			UserID:        ownerID,
			Title:         "Chart",
			Genre:         "Pop",
			Description:   "Description",
			PositionCount: 5,
		},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByID(context.Background(), chartID, &viewerID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected chart, got nil")
	}

	if !producer.sendReactionCalled {
		t.Errorf("expected SendReaction to be called")
	}
}
func TestGetChartByIDWithoutView_ReturnsChart(t *testing.T) {
	ownerID := uuid.New()
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	expectedChart := &chart.Chart{
		ID:            chartID,
		UserID:        ownerID,
		Title:         "Chart",
		Genre:         "Pop",
		Description:   "Description",
		PositionCount: 5,
	}

	repo := &fakeChartRepository{
		chartByID: expectedChart,
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByIDWithoutView(context.Background(), chartID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expectedChart {
		t.Errorf("expected chart %v, got %v", expectedChart, result)
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called")
	}
}
func TestGetChartByIDWithoutView_ReturnsErrorWhenChartNotFound(t *testing.T) {
	chartID := uuid.New()

	producer := &fakeEventProducer{}

	repo := &fakeChartRepository{
		chartByID: nil,
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, producer)

	result, err := service.GetChartByIDWithoutView(context.Background(), chartID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if producer.sendReactionCalled {
		t.Errorf("SendReaction should not be called")
	}
}
func TestGetMyChart_ReturnsFirstUserChart(t *testing.T) {
	userID := uuid.New()

	expectedChart := chart.Chart{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         "My Chart",
		Genre:         "Pop",
		Description:   "Description",
		PositionCount: 5,
	}

	repo := &fakeChartRepository{
		charts: []chart.Chart{expectedChart},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, &fakeEventProducer{})

	result, err := service.GetMyChart(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected chart, got nil")
	}

	if *result != expectedChart {
		t.Errorf("expected chart %v, got %v", expectedChart, *result)
	}
}
func TestGetMyChart_ReturnsErrorWhenUserHasNoChart(t *testing.T) {
	userID := uuid.New()

	repo := &fakeChartRepository{
		charts: []chart.Chart{},
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, &fakeEventProducer{})

	result, err := service.GetMyChart(context.Background(), userID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
func TestGetMyChart_ReturnsErrorWhenGetByUserIDFails(t *testing.T) {
	userID := uuid.New()
	expectedErr := errors.New("database error")

	repo := &fakeChartRepository{
		getByUserIDErr: expectedErr,
	}

	service := NewChartService(repo, nil, &fakeAnalyticsClient{}, &fakeEventProducer{})

	result, err := service.GetMyChart(context.Background(), userID)

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

func TestGetMostPopularCharts_ReturnsChartsInAnalyticsOrder(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	analyticsClient := &fakeAnalyticsClient{
		popularIDs: []uuid.UUID{id2, id3, id1},
	}

	repo := &fakeChartRepository{
		charts: []chart.Chart{
			{ID: id1, Title: "Chart 1", PositionCount: 5},
			{ID: id2, Title: "Chart 2", PositionCount: 5},
			{ID: id3, Title: "Chart 3", PositionCount: 5},
		},
	}

	service := NewChartService(repo, nil, analyticsClient, &fakeEventProducer{})

	result, err := service.GetMostPopularCharts(context.Background(), 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 charts, got %d", len(result))
	}

	if result[0].ID != id2 || result[1].ID != id3 || result[2].ID != id1 {
		t.Errorf("charts order does not match analytics IDs order")
	}
}

func TestGetMostPopularCharts_UsesDefaultLimitWhenLimitIsInvalid(t *testing.T) {
	analyticsClient := &fakeAnalyticsClient{
		popularIDs: []uuid.UUID{},
	}

	service := NewChartService(
		&fakeChartRepository{},
		nil,
		analyticsClient,
		&fakeEventProducer{},
	)

	_, err := service.GetMostPopularCharts(context.Background(), 0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if analyticsClient.limit != 25 {
		t.Errorf("expected default limit 25, got %d", analyticsClient.limit)
	}
}

func TestGetMostPopularCharts_ReturnsEmptyListWhenAnalyticsReturnsNoIDs(t *testing.T) {
	analyticsClient := &fakeAnalyticsClient{
		popularIDs: []uuid.UUID{},
	}

	service := NewChartService(
		&fakeChartRepository{},
		nil,
		analyticsClient,
		&fakeEventProducer{},
	)

	result, err := service.GetMostPopularCharts(context.Background(), 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d charts", len(result))
	}
}
func TestGetMostPopularCharts_ReturnsErrorWhenAnalyticsFails(t *testing.T) {
	expectedErr := errors.New("analytics error")

	analyticsClient := &fakeAnalyticsClient{
		err: expectedErr,
	}

	service := NewChartService(
		&fakeChartRepository{},
		nil,
		analyticsClient,
		&fakeEventProducer{},
	)

	result, err := service.GetMostPopularCharts(context.Background(), 10)

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
func TestGetUserLikedCharts_ReturnsChartsInAnalyticsOrder(t *testing.T) {
	userID := uuid.New()

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	analyticsClient := &fakeAnalyticsClient{
		likedIDs: []uuid.UUID{id2, id3, id1},
	}

	repo := &fakeChartRepository{
		charts: []chart.Chart{
			{ID: id1, Title: "Chart 1", PositionCount: 5},
			{ID: id2, Title: "Chart 2", PositionCount: 5},
			{ID: id3, Title: "Chart 3", PositionCount: 5},
		},
	}

	service := NewChartService(repo, nil, analyticsClient, &fakeEventProducer{})

	result, err := service.GetUserLikedCharts(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 charts, got %d", len(result))
	}

	if result[0].ID != id2 || result[1].ID != id3 || result[2].ID != id1 {
		t.Errorf("charts order does not match liked IDs order")
	}
}
func TestGetUserLikedCharts_ReturnsEmptyListWhenAnalyticsReturnsNoIDs(t *testing.T) {
	userID := uuid.New()

	analyticsClient := &fakeAnalyticsClient{
		likedIDs: []uuid.UUID{},
	}

	service := NewChartService(
		&fakeChartRepository{},
		nil,
		analyticsClient,
		&fakeEventProducer{},
	)

	result, err := service.GetUserLikedCharts(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d charts", len(result))
	}
}
func TestGetUserDislikedCharts_ReturnsChartsInAnalyticsOrder(t *testing.T) {
	userID := uuid.New()

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	analyticsClient := &fakeAnalyticsClient{
		dislikedIDs: []uuid.UUID{id3, id1, id2},
	}

	repo := &fakeChartRepository{
		charts: []chart.Chart{
			{ID: id1, Title: "Chart 1", PositionCount: 5},
			{ID: id2, Title: "Chart 2", PositionCount: 5},
			{ID: id3, Title: "Chart 3", PositionCount: 5},
		},
	}

	service := NewChartService(repo, nil, analyticsClient, &fakeEventProducer{})

	result, err := service.GetUserDislikedCharts(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 charts, got %d", len(result))
	}

	if result[0].ID != id3 || result[1].ID != id1 || result[2].ID != id2 {
		t.Errorf("charts order does not match disliked IDs order")
	}
}
func TestGetUserDislikedCharts_ReturnsEmptyListWhenAnalyticsReturnsNoIDs(t *testing.T) {
	userID := uuid.New()

	analyticsClient := &fakeAnalyticsClient{
		dislikedIDs: []uuid.UUID{},
	}

	service := NewChartService(
		&fakeChartRepository{},
		nil,
		analyticsClient,
		&fakeEventProducer{},
	)

	result, err := service.GetUserDislikedCharts(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d charts", len(result))
	}
}
