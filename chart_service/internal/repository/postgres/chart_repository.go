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

func (r *ChartRepository) GetIDsByGenre(
	ctx context.Context,
	genre string,
	limit int,
) ([]uuid.UUID, error) {

	query := `
	SELECT id
	FROM charts
	WHERE genre = $1
	LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, genre, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *ChartRepository) GetGenresByChartIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (map[uuid.UUID]string, error) {

	query := `
	SELECT id, genre
	FROM charts
	WHERE id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]string)

	for rows.Next() {
		var id uuid.UUID
		var genre string

		if err := rows.Scan(&id, &genre); err != nil {
			return nil, err
		}

		result[id] = genre
	}

	return result, nil
}
