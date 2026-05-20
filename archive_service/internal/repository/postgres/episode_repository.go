package postgres

import (
	"context"
	"time"

	"charts-archive-service/internal/domain/episode"
	"charts-archive-service/internal/domain/track_episode"
	"charts-archive-service/internal/repository/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EpisodeRepository struct {
	db *pgxpool.Pool
}

func NewEpisodeRepository(db *pgxpool.Pool) *EpisodeRepository {
	return &EpisodeRepository{db: db}
}

func (r *EpisodeRepository) Create(ctx context.Context, e *episode.Episode) (*episode.Episode, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	episodeQuery := `
	INSERT INTO episodes (id, chart_id, created_at)
	VALUES ($1, $2, $3)
	`

	_, err = tx.Exec(ctx, episodeQuery,
		e.ID,
		e.ChartID,
		e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	trackQuery := `
	INSERT INTO track_episode (
		id,
		episode_id,
		track_id,
		current_position,
		previous_position,
		highest_position,
		episodes_count,
		times_at_peak_position
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, t := range e.Tracks {
		_, err = tx.Exec(ctx, trackQuery,
			t.ID,
			e.ID,
			t.TrackID,
			t.CurrentPosition,
			t.PreviousPosition,
			t.HighestPosition,
			t.EpisodesCount,
			t.TimesAtPeakPosition,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EpisodeRepository) GetByID(ctx context.Context, id uuid.UUID) (*episode.Episode, error) {
	query := `
	SELECT id, chart_id, created_at
	FROM episodes
	WHERE id = $1
	`

	e := &episode.Episode{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID,
		&e.ChartID,
		&e.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EpisodeRepository) GetByChartID(ctx context.Context, chartID uuid.UUID) ([]episode.Episode, error) {
	query := `
	SELECT id, chart_id, created_at
	FROM episodes
	WHERE chart_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, chartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodesList []episode.Episode

	for rows.Next() {
		var e episode.Episode

		err := rows.Scan(
			&e.ID,
			&e.ChartID,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		episodesList = append(episodesList, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return episodesList, nil
}

func (r *EpisodeRepository) GetLatestByChartID(ctx context.Context, chartID uuid.UUID) (*episode.Episode, error) {
	query := `
	SELECT id, chart_id, created_at
	FROM episodes
	WHERE chart_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	e := &episode.Episode{}

	err := r.db.QueryRow(ctx, query, chartID).Scan(
		&e.ID,
		&e.ChartID,
		&e.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	tracksQuery := `
	SELECT id, track_id, current_position, previous_position,
	       highest_position, episodes_count, times_at_peak_position
	FROM track_episode
	WHERE episode_id = $1
	ORDER BY current_position
	`

	rows, err := r.db.Query(ctx, tracksQuery, e.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []track_episode.TrackEpisode

	for rows.Next() {
		var t track_episode.TrackEpisode

		err := rows.Scan(
			&t.ID,
			&t.TrackID,
			&t.CurrentPosition,
			&t.PreviousPosition,
			&t.HighestPosition,
			&t.EpisodesCount,
			&t.TimesAtPeakPosition,
		)
		if err != nil {
			return nil, err
		}

		t.EpisodeID = e.ID
		tracks = append(tracks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	e.Tracks = tracks

	return e, nil
}

func (r *EpisodeRepository) GetLatestWithTracksByLimit(
	ctx context.Context,
	limit int,
) ([]dto.EpisodeResponse, error) {

	query := `
	SELECT 
		e.id, e.chart_id, e.created_at,
		te.track_id,
		t.artist,
		t.title,
		te.current_position,
		te.previous_position,
		te.highest_position,
		te.episodes_count,
		te.times_at_peak_position
	FROM episodes e
	JOIN track_episode te ON te.episode_id = e.id
	JOIN tracks t ON t.id = te.track_id
	WHERE e.id IN (
		SELECT id FROM episodes
		ORDER BY created_at DESC
		LIMIT $1
	)
	ORDER BY e.created_at DESC, te.current_position
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	episodesMap := make(map[uuid.UUID]*dto.EpisodeResponse)
	order := make([]uuid.UUID, 0)

	for rows.Next() {
		var (
			epID, chartID uuid.UUID
			createdAt     time.Time

			track dto.TrackEpisodeResponse
		)

		err := rows.Scan(
			&epID,
			&chartID,
			&createdAt,
			&track.TrackID,
			&track.Artist,
			&track.Title,
			&track.CurrentPosition,
			&track.PreviousPosition,
			&track.HighestPosition,
			&track.EpisodesCount,
			&track.TimesAtPeakPosition,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := episodesMap[epID]; !ok {
			episodesMap[epID] = &dto.EpisodeResponse{
				ID:        epID,
				ChartID:   chartID,
				CreatedAt: createdAt,
				Tracks:    []dto.TrackEpisodeResponse{},
			}
			order = append(order, epID)
		}

		episodesMap[epID].Tracks = append(
			episodesMap[epID].Tracks,
			track,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]dto.EpisodeResponse, 0, len(order))
	for _, id := range order {
		result = append(result, *episodesMap[id])
	}

	return result, nil
}

func (r *EpisodeRepository) GetTracksWithMetaByEpisodeID(
	ctx context.Context,
	episodeID uuid.UUID,
) ([]dto.TrackEpisodeResponse, error) {

	query := `
	SELECT 
		te.track_id,
		t.artist,
		t.title,
		te.current_position,
		te.previous_position,
		te.highest_position,
		te.episodes_count,
		te.times_at_peak_position
	FROM track_episode te
	JOIN tracks t ON te.track_id = t.id
	WHERE te.episode_id = $1
	ORDER BY te.current_position
	`

	rows, err := r.db.Query(ctx, query, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.TrackEpisodeResponse

	for rows.Next() {
		var t dto.TrackEpisodeResponse

		err := rows.Scan(
			&t.TrackID,
			&t.Artist,
			&t.Title,
			&t.CurrentPosition,
			&t.PreviousPosition,
			&t.HighestPosition,
			&t.EpisodesCount,
			&t.TimesAtPeakPosition,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EpisodeRepository) GetNearestLeftWithTracks(
	ctx context.Context,
	chartID uuid.UUID,
	date time.Time,
) (*dto.EpisodeResponse, error) {

	query := `
	SELECT 
		e.id, e.chart_id, e.created_at,
		te.track_id,
		t.artist,
		t.title,
		te.current_position,
		te.previous_position,
		te.highest_position,
		te.episodes_count,
		te.times_at_peak_position
	FROM episodes e
	JOIN track_episode te ON te.episode_id = e.id
	JOIN tracks t ON t.id = te.track_id
	WHERE e.chart_id = $1
	  AND e.created_at < $2
	ORDER BY e.created_at DESC, te.current_position
	`

	rows, err := r.db.Query(ctx, query, chartID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *dto.EpisodeResponse

	for rows.Next() {
		var (
			epID, chartID uuid.UUID
			createdAt     time.Time

			track dto.TrackEpisodeResponse
		)

		err := rows.Scan(
			&epID,
			&chartID,
			&createdAt,
			&track.TrackID,
			&track.Artist,
			&track.Title,
			&track.CurrentPosition,
			&track.PreviousPosition,
			&track.HighestPosition,
			&track.EpisodesCount,
			&track.TimesAtPeakPosition,
		)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = &dto.EpisodeResponse{
				ID:        epID,
				ChartID:   chartID,
				CreatedAt: createdAt,
				Tracks:    []dto.TrackEpisodeResponse{},
			}
		}

		if result.ID != epID {
			break
		}

		result.Tracks = append(result.Tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *EpisodeRepository) GetLatestEpisodesPage(
	ctx context.Context,
	limit int,
	offset int,
) ([]episode.Episode, error) {

	query := `
	SELECT id, chart_id, created_at
	FROM episodes
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []episode.Episode

	for rows.Next() {
		var e episode.Episode

		err := rows.Scan(
			&e.ID,
			&e.ChartID,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
