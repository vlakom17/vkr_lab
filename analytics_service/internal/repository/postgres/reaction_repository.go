package postgres

import (
	"context"

	"charts-analytics-service/internal/domain/reaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReactionRepository struct {
	db *pgxpool.Pool
}

func NewReactionRepository(db *pgxpool.Pool) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) Upsert(
	ctx context.Context,
	rct *reaction.Reaction,
) (*reaction.Reaction, error) {

	query := `
		INSERT INTO reactions (user_id, chart_id, type, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, chart_id)
		DO UPDATE SET
			type = EXCLUDED.type,
			created_at = EXCLUDED.created_at
		RETURNING user_id, chart_id, type, created_at;
	`

	var res reaction.Reaction

	err := r.db.QueryRow(
		ctx,
		query,
		rct.UserID,
		rct.ChartID,
		rct.Type,
		rct.CreatedAt,
	).Scan(
		&res.UserID,
		&res.ChartID,
		&res.Type,
		&res.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *ReactionRepository) GetByUserAndChart(
	ctx context.Context,
	userID, chartID uuid.UUID,
) (*reaction.Reaction, error) {

	query := `
		SELECT user_id, chart_id, type, created_at
		FROM reactions
		WHERE user_id = $1 AND chart_id = $2;
	`

	var res reaction.Reaction

	err := r.db.QueryRow(ctx, query, userID, chartID).Scan(
		&res.UserID,
		&res.ChartID,
		&res.Type,
		&res.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func (r *ReactionRepository) CountByChart(
	ctx context.Context,
	chartID uuid.UUID,
) (likes, dislikes, views int, err error) {

	query := `
		SELECT
			COUNT(*) FILTER (WHERE type = 'like') AS likes,
			COUNT(*) FILTER (WHERE type = 'dislike') AS dislikes,
			COUNT(*) AS views
		FROM reactions
		WHERE chart_id = $1;
	`

	err = r.db.QueryRow(ctx, query, chartID).Scan(
		&likes,
		&dislikes,
		&views,
	)

	return
}

func (r *ReactionRepository) GetMostPopularChartIDs(
	ctx context.Context,
	limit int,
) ([]uuid.UUID, error) {

	query := `
		SELECT chart_id
		FROM reactions
		GROUP BY chart_id
		ORDER BY
			COUNT(*) FILTER (WHERE type = 'like')
		  - COUNT(*) FILTER (WHERE type = 'dislike') DESC
		LIMIT $1;
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]uuid.UUID, 0, limit)

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, rows.Err()
}

func (r *ReactionRepository) GetUserChartIDsByType(
	ctx context.Context,
	userID uuid.UUID,
	t reaction.ReactionType,
) ([]uuid.UUID, error) {

	query := `
		SELECT chart_id
		FROM reactions
		WHERE user_id = $1 AND type = $2;
	`

	rows, err := r.db.Query(ctx, query, userID, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []uuid.UUID

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, rows.Err()
}

func (r *ReactionRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]reaction.Reaction, error) {

	query := `
		SELECT user_id, chart_id, type, created_at
		FROM reactions
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []reaction.Reaction

	for rows.Next() {
		var rct reaction.Reaction
		if err := rows.Scan(
			&rct.UserID,
			&rct.ChartID,
			&rct.Type,
			&rct.CreatedAt,
		); err != nil {
			return nil, err
		}

		reactions = append(reactions, rct)
	}

	return reactions, rows.Err()
}
