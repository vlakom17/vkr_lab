package service

import (
	"context"

	"charts-analytics-service/internal/domain/event"
	"charts-analytics-service/internal/domain/reaction"

	"github.com/google/uuid"
)

type UserClient interface {
	GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
}

type ReactionService struct {
	repo       reaction.ReactionRepository
	userClient UserClient
}

func NewReactionService(
	repo reaction.ReactionRepository,
	userClient UserClient,
) *ReactionService {
	return &ReactionService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *ReactionService) GetReactionStats(
	ctx context.Context,
	chartID uuid.UUID,
) (*reaction.ReactionStats, error) {

	likes, dislikes, views, err := s.repo.CountByChart(ctx, chartID)
	if err != nil {
		return nil, err
	}

	stats := &reaction.ReactionStats{
		ChartID:       chartID,
		LikesCount:    likes,
		DislikesCount: dislikes,
		ViewsCount:    views,
	}

	return stats, nil
}

func (s *ReactionService) GetMostPopularChartIDs(
	ctx context.Context,
	limit int,
) ([]uuid.UUID, error) {

	return s.repo.GetMostPopularChartIDs(ctx, limit)
}

func (s *ReactionService) GetUserLikedChartIDs(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {

	return s.repo.GetUserChartIDsByType(
		ctx,
		userID,
		reaction.ReactionLike,
	)
}

func (s *ReactionService) GetUserDislikedChartIDs(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {

	return s.repo.GetUserChartIDsByType(
		ctx,
		userID,
		reaction.ReactionDislike,
	)
}

func (s *ReactionService) HandleReactionEvent(
	ctx context.Context,
	e event.ReactionEvent,
) error {

	if e.UserID == uuid.Nil || e.ChartID == uuid.Nil {
		return nil
	}

	existing, err := s.repo.GetByUserAndChart(ctx, e.UserID, e.ChartID)
	if err != nil {
		return err
	}

	switch e.Type {

	case "like":
		rct := &reaction.Reaction{
			UserID:    e.UserID,
			ChartID:   e.ChartID,
			Type:      reaction.ReactionLike,
			CreatedAt: e.CreatedAt,
		}
		_, err = s.repo.Upsert(ctx, rct)
		return err

	case "dislike":
		rct := &reaction.Reaction{
			UserID:    e.UserID,
			ChartID:   e.ChartID,
			Type:      reaction.ReactionDislike,
			CreatedAt: e.CreatedAt,
		}
		_, err = s.repo.Upsert(ctx, rct)
		return err

	case "view":
		// view не затирает оценку
		if existing != nil &&
			(existing.Type == reaction.ReactionLike || existing.Type == reaction.ReactionDislike) {
			return nil
		}

		rct := &reaction.Reaction{
			UserID:    e.UserID,
			ChartID:   e.ChartID,
			Type:      reaction.ReactionView,
			CreatedAt: e.CreatedAt,
		}
		_, err = s.repo.Upsert(ctx, rct)
		return err

	case "remove":
		if existing == nil {
			return nil
		}

		if existing.Type == reaction.ReactionView {
			return nil
		}

		rct := &reaction.Reaction{
			UserID:    e.UserID,
			ChartID:   e.ChartID,
			Type:      reaction.ReactionView,
			CreatedAt: e.CreatedAt,
		}
		_, err = s.repo.Upsert(ctx, rct)
		return err

	default:
		return nil
	}
}

func (s *ReactionService) GetMyReactionOnChart(
	ctx context.Context,
	token string,
	chartID uuid.UUID,
) (*reaction.Reaction, error) {

	userID, err := s.userClient.GetUserIDByToken(
		ctx,
		token,
	)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByUserAndChart(
		ctx,
		userID,
		chartID,
	)
}
