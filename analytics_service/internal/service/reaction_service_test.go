package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"charts-analytics-service/internal/domain/event"
	"charts-analytics-service/internal/domain/reaction"

	"github.com/google/uuid"
)

type fakeReactionRepository struct {
	existingReaction *reaction.Reaction
	upsertedReaction *reaction.Reaction

	getByUserAndChartErr error
	upsertErr            error

	upsertCalled  bool
	countChartID  uuid.UUID
	likesCount    int
	dislikesCount int
	viewsCount    int
	countErr      error

	popularLimit int
	popularIDs   []uuid.UUID
	popularErr   error

	requestedUserID uuid.UUID
	requestedType   reaction.ReactionType
	userChartIDs    []uuid.UUID
	userChartIDsErr error
}

func (r *fakeReactionRepository) Upsert(
	ctx context.Context,
	rct *reaction.Reaction,
) (*reaction.Reaction, error) {
	r.upsertCalled = true

	if r.upsertErr != nil {
		return nil, r.upsertErr
	}

	r.upsertedReaction = rct
	r.existingReaction = rct

	return rct, nil
}

func (r *fakeReactionRepository) GetByUserAndChart(
	ctx context.Context,
	userID uuid.UUID,
	chartID uuid.UUID,
) (*reaction.Reaction, error) {
	if r.getByUserAndChartErr != nil {
		return nil, r.getByUserAndChartErr
	}

	return r.existingReaction, nil
}

func (r *fakeReactionRepository) CountByChart(
	ctx context.Context,
	chartID uuid.UUID,
) (likes, dislikes, views int, err error) {
	r.countChartID = chartID

	if r.countErr != nil {
		return 0, 0, 0, r.countErr
	}

	return r.likesCount, r.dislikesCount, r.viewsCount, nil
}

func (r *fakeReactionRepository) GetMostPopularChartIDs(
	ctx context.Context,
	limit int,
) ([]uuid.UUID, error) {
	r.popularLimit = limit

	if r.popularErr != nil {
		return nil, r.popularErr
	}

	return r.popularIDs, nil
}

func (r *fakeReactionRepository) GetUserChartIDsByType(
	ctx context.Context,
	userID uuid.UUID,
	t reaction.ReactionType,
) ([]uuid.UUID, error) {
	r.requestedUserID = userID
	r.requestedType = t

	if r.userChartIDsErr != nil {
		return nil, r.userChartIDsErr
	}

	return r.userChartIDs, nil
}

func (r *fakeReactionRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]reaction.Reaction, error) {
	return nil, nil
}

type fakeReactionUserClient struct {
	userID uuid.UUID
	err    error
}

func (c *fakeReactionUserClient) GetUserIDByToken(
	ctx context.Context,
	token string,
) (uuid.UUID, error) {
	if c.err != nil {
		return uuid.Nil, c.err
	}

	return c.userID, nil
}
func TestHandleReactionEvent_UpsertsLikeReaction(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()
	createdAt := time.Now()

	repo := &fakeReactionRepository{}
	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "like",
		CreatedAt: createdAt,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.upsertCalled {
		t.Fatalf("expected Upsert to be called")
	}

	if repo.upsertedReaction == nil {
		t.Fatalf("expected upserted reaction")
	}

	if repo.upsertedReaction.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, repo.upsertedReaction.UserID)
	}

	if repo.upsertedReaction.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, repo.upsertedReaction.ChartID)
	}

	if repo.upsertedReaction.Type != reaction.ReactionLike {
		t.Errorf("expected reaction type %s, got %s", reaction.ReactionLike, repo.upsertedReaction.Type)
	}

	if !repo.upsertedReaction.CreatedAt.Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, repo.upsertedReaction.CreatedAt)
	}
}

func TestHandleReactionEvent_UpsertsDislikeReaction(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()
	createdAt := time.Now()

	repo := &fakeReactionRepository{}
	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "dislike",
		CreatedAt: createdAt,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.upsertCalled {
		t.Fatalf("expected Upsert to be called")
	}

	if repo.upsertedReaction.Type != reaction.ReactionDislike {
		t.Errorf("expected reaction type %s, got %s", reaction.ReactionDislike, repo.upsertedReaction.Type)
	}
}

func TestHandleReactionEvent_UpsertsViewWhenNoExistingReaction(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()
	createdAt := time.Now()

	repo := &fakeReactionRepository{}
	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "view",
		CreatedAt: createdAt,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.upsertCalled {
		t.Fatalf("expected Upsert to be called")
	}

	if repo.upsertedReaction.Type != reaction.ReactionView {
		t.Errorf("expected reaction type %s, got %s", reaction.ReactionView, repo.upsertedReaction.Type)
	}
}

func TestHandleReactionEvent_DoesNotOverwriteLikeWithView(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeReactionRepository{
		existingReaction: &reaction.Reaction{
			UserID:  userID,
			ChartID: chartID,
			Type:    reaction.ReactionLike,
		},
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "view",
		CreatedAt: time.Now(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.upsertCalled {
		t.Errorf("Upsert should not be called when existing reaction is like")
	}
}
func TestHandleReactionEvent_RemoveDoesNothingWhenNoExistingReaction(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeReactionRepository{}
	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "remove",
		CreatedAt: time.Now(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.upsertCalled {
		t.Errorf("Upsert should not be called when there is no existing reaction")
	}
}
func TestHandleReactionEvent_RemoveChangesLikeToView(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()
	createdAt := time.Now()

	repo := &fakeReactionRepository{
		existingReaction: &reaction.Reaction{
			UserID:  userID,
			ChartID: chartID,
			Type:    reaction.ReactionLike,
		},
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "remove",
		CreatedAt: createdAt,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.upsertCalled {
		t.Fatalf("expected Upsert to be called")
	}

	if repo.upsertedReaction.Type != reaction.ReactionView {
		t.Errorf("expected reaction type %s, got %s", reaction.ReactionView, repo.upsertedReaction.Type)
	}

	if !repo.upsertedReaction.CreatedAt.Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, repo.upsertedReaction.CreatedAt)
	}
}
func TestHandleReactionEvent_RemoveDoesNothingWhenExistingReactionIsView(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeReactionRepository{
		existingReaction: &reaction.Reaction{
			UserID:  userID,
			ChartID: chartID,
			Type:    reaction.ReactionView,
		},
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "remove",
		CreatedAt: time.Now(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.upsertCalled {
		t.Errorf("Upsert should not be called when existing reaction is view")
	}
}
func TestHandleReactionEvent_IgnoresEventWhenUserIDIsNil(t *testing.T) {
	chartID := uuid.New()

	repo := &fakeReactionRepository{}
	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    uuid.Nil,
		ChartID:   chartID,
		Type:      "like",
		CreatedAt: time.Now(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.upsertCalled {
		t.Errorf("Upsert should not be called when userID is nil")
	}
}
func TestHandleReactionEvent_IgnoresUnknownReactionType(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeReactionRepository{}
	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "unknown",
		CreatedAt: time.Now(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.upsertCalled {
		t.Errorf("Upsert should not be called for unknown reaction type")
	}
}
func TestHandleReactionEvent_ReturnsErrorWhenGetByUserAndChartFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("database error")

	repo := &fakeReactionRepository{
		getByUserAndChartErr: expectedErr,
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "like",
		CreatedAt: time.Now(),
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if repo.upsertCalled {
		t.Errorf("Upsert should not be called when GetByUserAndChart fails")
	}
}
func TestHandleReactionEvent_ReturnsErrorWhenUpsertFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("upsert error")

	repo := &fakeReactionRepository{
		upsertErr: expectedErr,
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	err := service.HandleReactionEvent(context.Background(), event.ReactionEvent{
		UserID:    userID,
		ChartID:   chartID,
		Type:      "like",
		CreatedAt: time.Now(),
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if !repo.upsertCalled {
		t.Errorf("expected Upsert to be called")
	}
}
func TestGetReactionStats_ReturnsStats(t *testing.T) {
	chartID := uuid.New()

	repo := &fakeReactionRepository{
		likesCount:    10,
		dislikesCount: 3,
		viewsCount:    100,
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	result, err := service.GetReactionStats(context.Background(), chartID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, result.ChartID)
	}

	if result.LikesCount != 10 {
		t.Errorf("expected likes 10, got %d", result.LikesCount)
	}

	if result.DislikesCount != 3 {
		t.Errorf("expected dislikes 3, got %d", result.DislikesCount)
	}

	if result.ViewsCount != 100 {
		t.Errorf("expected views 100, got %d", result.ViewsCount)
	}
}
func TestGetMostPopularChartIDs_ReturnsIDs(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	repo := &fakeReactionRepository{
		popularIDs: []uuid.UUID{id1, id2},
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	result, err := service.GetMostPopularChartIDs(context.Background(), 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.popularLimit != 10 {
		t.Errorf("expected limit 10, got %d", repo.popularLimit)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(result))
	}

	if result[0] != id1 || result[1] != id2 {
		t.Errorf("unexpected ids order")
	}
}
func TestGetUserLikedChartIDs_RequestsLikedType(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeReactionRepository{
		userChartIDs: []uuid.UUID{chartID},
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	result, err := service.GetUserLikedChartIDs(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.requestedUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, repo.requestedUserID)
	}

	if repo.requestedType != reaction.ReactionLike {
		t.Errorf("expected type %s, got %s", reaction.ReactionLike, repo.requestedType)
	}

	if len(result) != 1 || result[0] != chartID {
		t.Errorf("unexpected result")
	}
}
func TestGetUserDislikedChartIDs_RequestsDislikedType(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	repo := &fakeReactionRepository{
		userChartIDs: []uuid.UUID{chartID},
	}

	service := NewReactionService(repo, &fakeReactionUserClient{})

	result, err := service.GetUserDislikedChartIDs(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.requestedUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, repo.requestedUserID)
	}

	if repo.requestedType != reaction.ReactionDislike {
		t.Errorf("expected type %s, got %s", reaction.ReactionDislike, repo.requestedType)
	}

	if len(result) != 1 || result[0] != chartID {
		t.Errorf("unexpected result")
	}
}
func TestGetMyReactionOnChart_ReturnsReaction(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()

	expectedReaction := &reaction.Reaction{
		UserID:  userID,
		ChartID: chartID,
		Type:    reaction.ReactionLike,
	}

	repo := &fakeReactionRepository{
		existingReaction: expectedReaction,
	}

	userClient := &fakeReactionUserClient{
		userID: userID,
	}

	service := NewReactionService(repo, userClient)

	result, err := service.GetMyReactionOnChart(
		context.Background(),
		"token",
		chartID,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expectedReaction {
		t.Errorf("expected reaction %v, got %v", expectedReaction, result)
	}
}
func TestGetMyReactionOnChart_ReturnsErrorWhenUserClientFails(t *testing.T) {
	chartID := uuid.New()
	expectedErr := errors.New("invalid token")

	userClient := &fakeReactionUserClient{
		err: expectedErr,
	}

	service := NewReactionService(&fakeReactionRepository{}, userClient)

	result, err := service.GetMyReactionOnChart(
		context.Background(),
		"bad-token",
		chartID,
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
func TestGetMyReactionOnChart_ReturnsErrorWhenRepositoryFails(t *testing.T) {
	userID := uuid.New()
	chartID := uuid.New()
	expectedErr := errors.New("database error")

	repo := &fakeReactionRepository{
		getByUserAndChartErr: expectedErr,
	}

	userClient := &fakeReactionUserClient{
		userID: userID,
	}

	service := NewReactionService(repo, userClient)

	result, err := service.GetMyReactionOnChart(
		context.Background(),
		"token",
		chartID,
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
