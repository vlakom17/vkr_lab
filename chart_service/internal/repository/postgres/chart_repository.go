package postgres

import (
	"context"

	"charts-chart-service/internal/domain/chart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChartRepository struct {
	db *pgxpool.Pool
}

func NewChartRepository(db *pgxpool.Pool) *ChartRepository {
	return &ChartRepository{db: db}
}

func (r *ChartRepository) Create(ctx context.Context, c *chart.Chart) (*chart.Chart, error) {
	query := `
	INSERT INTO charts (id, user_id, title, genre, description, position_count)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		c.ID,
		c.UserID,
		c.Title,
		c.Genre,
		c.Description,
		c.PositionCount,
	).Scan(&c.CreatedAt)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *ChartRepository) Update(ctx context.Context, c *chart.Chart) (*chart.Chart, error) {
	query := `
	UPDATE charts
	SET title = $2,
	    genre = $3,
	    description = $4,
	    position_count = $5
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		c.ID,
		c.Title,
		c.Genre,
		c.Description,
		c.PositionCount,
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *ChartRepository) GetByID(ctx context.Context, id uuid.UUID) (*chart.Chart, error) {
	query := `
	SELECT id, user_id, title, genre, description, position_count, created_at
	FROM charts
	WHERE id = $1
	`

	c := &chart.Chart{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.UserID,
		&c.Title,
		&c.Genre,
		&c.Description,
		&c.PositionCount,
		&c.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *ChartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]chart.Chart, error) {
	query := `
	SELECT id, user_id, title, genre, description, position_count, created_at
	FROM charts
	WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charts []chart.Chart

	for rows.Next() {
		var c chart.Chart

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Title,
			&c.Genre,
			&c.Description,
			&c.PositionCount,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		charts = append(charts, c)
	}

	return charts, nil
}

func (r *ChartRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]chart.Chart, error) {
	query := `
	SELECT id, user_id, title, genre, description, position_count, created_at
	FROM charts
	WHERE id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charts []chart.Chart

	for rows.Next() {
		var c chart.Chart

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Title,
			&c.Genre,
			&c.Description,
			&c.PositionCount,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		charts = append(charts, c)
	}

	return charts, nil
}
