package postgres

import (
	"context"

	"charts-archive-service/internal/domain/track"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackRepository struct {
	db *pgxpool.Pool
}

func NewTrackRepository(db *pgxpool.Pool) *TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) FindOrCreate(
	ctx context.Context,
	artist, title, normalizedKey string,
) (*track.Track, error) {

	query := `
	INSERT INTO tracks (id, artist, title, normalized_key)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (normalized_key) DO NOTHING
	RETURNING id, artist, title, normalized_key
	`

	id := uuid.New()

	t := &track.Track{}

	err := r.db.QueryRow(ctx, query,
		id,
		artist,
		title,
		normalizedKey,
	).Scan(
		&t.ID,
		&t.Artist,
		&t.Title,
		&t.NormalizedKey,
	)

	if err == pgx.ErrNoRows {

		query = `
		SELECT id, artist, title, normalized_key
		FROM tracks
		WHERE normalized_key = $1
		`

		err = r.db.QueryRow(ctx, query, normalizedKey).Scan(
			&t.ID,
			&t.Artist,
			&t.Title,
			&t.NormalizedKey,
		)
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TrackRepository) Search(ctx context.Context, query string) ([]track.Track, error) {
	sqlQuery := `
	SELECT id, artist, title, normalized_key
	FROM tracks
	WHERE artist ILIKE $1 || '%'
	   OR title ILIKE $1 || '%'
	ORDER BY artist, title
	LIMIT 10
	`

	rows, err := r.db.Query(ctx, sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracksResult []track.Track

	for rows.Next() {
		var t track.Track

		err := rows.Scan(
			&t.ID,
			&t.Artist,
			&t.Title,
			&t.NormalizedKey,
		)
		if err != nil {
			return nil, err
		}

		tracksResult = append(tracksResult, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracksResult, nil
}
