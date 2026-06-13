package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"charts-archive-service/internal/domain/episode"
	"charts-archive-service/internal/domain/event"
	"charts-archive-service/internal/domain/track"
	"charts-archive-service/internal/domain/track_episode"
	"charts-archive-service/internal/repository/dto"

	"github.com/google/uuid"
)

type fakeEpisodeRepository struct {
	latestEpisode  *episode.Episode
	createdEpisode *episode.Episode

	getLatestByChartIDErr  error
	createErr              error
	latestPageLimit        int
	latestPageOffset       int
	latestPageResult       []episode.Episode
	latestPageErr          error
	episodeByID            *episode.Episode
	getByIDErr             error
	tracksWithMeta         []dto.TrackEpisodeResponse
	tracksWithMetaErr      error
	episodesByChart        []episode.Episode
	getByChartIDErr        error
	requestedChartID       uuid.UUID
	latestWithTracksLimit  int
	latestWithTracksResult []dto.EpisodeResponse
	latestWithTracksErr    error
	nearestLeftChartID     uuid.UUID
	nearestLeftDate        time.Time
	nearestLeftResult      *dto.EpisodeResponse
	nearestLeftErr         error
}

func (r *fakeEpisodeRepository) Create(ctx context.Context, ep *episode.Episode) (*episode.Episode, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}

	r.createdEpisode = ep
	return ep, nil
}

func (r *fakeEpisodeRepository) GetByID(ctx context.Context, id uuid.UUID) (*episode.Episode, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}

	return r.episodeByID, nil
}

func (r *fakeEpisodeRepository) GetByChartID(ctx context.Context, chartID uuid.UUID) ([]episode.Episode, error) {
	r.requestedChartID = chartID

	if r.getByChartIDErr != nil {
		return nil, r.getByChartIDErr
	}

	return r.episodesByChart, nil
}
func (r *fakeEpisodeRepository) GetLatestByChartID(ctx context.Context, chartID uuid.UUID) (*episode.Episode, error) {
	if r.getLatestByChartIDErr != nil {
		return nil, r.getLatestByChartIDErr
	}

	return r.latestEpisode, nil
}

func (r *fakeEpisodeRepository) GetLatestWithTracksByLimit(
	ctx context.Context,
	limit int,
) ([]dto.EpisodeResponse, error) {
	r.latestWithTracksLimit = limit

	if r.latestWithTracksErr != nil {
		return nil, r.latestWithTracksErr
	}

	return r.latestWithTracksResult, nil
}

func (r *fakeEpisodeRepository) GetTracksWithMetaByEpisodeID(
	ctx context.Context,
	episodeID uuid.UUID,
) ([]dto.TrackEpisodeResponse, error) {
	if r.tracksWithMetaErr != nil {
		return nil, r.tracksWithMetaErr
	}

	return r.tracksWithMeta, nil
}

func (r *fakeEpisodeRepository) GetNearestLeftWithTracks(
	ctx context.Context,
	chartID uuid.UUID,
	date time.Time,
) (*dto.EpisodeResponse, error) {
	r.nearestLeftChartID = chartID
	r.nearestLeftDate = date

	if r.nearestLeftErr != nil {
		return nil, r.nearestLeftErr
	}

	return r.nearestLeftResult, nil
}

func (r *fakeEpisodeRepository) GetLatestEpisodesPage(
	ctx context.Context,
	limit int,
	offset int,
) ([]episode.Episode, error) {
	r.latestPageLimit = limit
	r.latestPageOffset = offset

	if r.latestPageErr != nil {
		return nil, r.latestPageErr
	}

	return r.latestPageResult, nil
}

type fakeEpisodeTrackRepository struct {
	tracksByKey     map[string]*track.Track
	findOrCreateErr error
}

func (r *fakeEpisodeTrackRepository) FindOrCreate(
	ctx context.Context,
	artist string,
	title string,
	normalizedKey string,
) (*track.Track, error) {
	if r.findOrCreateErr != nil {
		return nil, r.findOrCreateErr
	}

	if r.tracksByKey == nil {
		r.tracksByKey = make(map[string]*track.Track)
	}

	if existing, ok := r.tracksByKey[normalizedKey]; ok {
		return existing, nil
	}

	newTrack := &track.Track{
		ID:            uuid.New(),
		Artist:        artist,
		Title:         title,
		NormalizedKey: normalizedKey,
	}

	r.tracksByKey[normalizedKey] = newTrack
	return newTrack, nil
}

func (r *fakeEpisodeTrackRepository) Search(ctx context.Context, query string) ([]track.Track, error) {
	return nil, nil
}

func validTracks() []event.TrackSnapshot {
	return []event.TrackSnapshot{
		{Artist: "Artist 1", Title: "Track 1", CurrentPosition: 1},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}
}
func TestHandleEpisodeCreatedEvent_CreatesFirstEpisodeWithInitialMetrics(t *testing.T) {
	chartID := uuid.New()
	createdAt := time.Now()

	episodeRepo := &fakeEpisodeRepository{}
	trackRepo := &fakeEpisodeTrackRepository{}

	service := NewEpisodeService(episodeRepo, trackRepo)

	tracks := validTracks()

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: createdAt,
		Tracks:    tracks,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if episodeRepo.createdEpisode == nil {
		t.Fatalf("expected episode to be created")
	}

	if episodeRepo.createdEpisode.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, episodeRepo.createdEpisode.ChartID)
	}

	if len(episodeRepo.createdEpisode.Tracks) != 5 {
		t.Fatalf("expected 5 tracks, got %d", len(episodeRepo.createdEpisode.Tracks))
	}

	for _, te := range episodeRepo.createdEpisode.Tracks {
		if te.PreviousPosition != 0 {
			t.Errorf("expected previous position 0, got %d", te.PreviousPosition)
		}

		if te.HighestPosition != te.CurrentPosition {
			t.Errorf("expected highest position %d, got %d", te.CurrentPosition, te.HighestPosition)
		}

		if te.EpisodesCount != 1 {
			t.Errorf("expected episodes count 1, got %d", te.EpisodesCount)
		}

		if te.TimesAtPeakPosition != 1 {
			t.Errorf("expected times at peak 1, got %d", te.TimesAtPeakPosition)
		}
	}
}

func TestHandleEpisodeCreatedEvent_UpdatesHighestPositionWhenTrackGetsNewPeak(t *testing.T) {
	chartID := uuid.New()
	trackID := uuid.New()

	episodeRepo := &fakeEpisodeRepository{
		latestEpisode: &episode.Episode{
			ID:      uuid.New(),
			ChartID: chartID,
			Tracks: []track_episode.TrackEpisode{
				{
					TrackID:             trackID,
					CurrentPosition:     4,
					HighestPosition:     3,
					TimesAtPeakPosition: 2,
					EpisodesCount:       5,
				},
			},
		},
	}

	normalizedKey := "artist|song"

	trackRepo := &fakeEpisodeTrackRepository{
		tracksByKey: map[string]*track.Track{
			normalizedKey: {
				ID:            trackID,
				Artist:        "artist",
				Title:         "song",
				NormalizedKey: normalizedKey,
			},
		},
	}

	service := NewEpisodeService(episodeRepo, trackRepo)

	tracks := []event.TrackSnapshot{
		{Artist: "Artist", Title: "Song", CurrentPosition: 2},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 3},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 4},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 5},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 1},
	}

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    tracks,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultTrack := episodeRepo.createdEpisode.Tracks[0]

	if resultTrack.PreviousPosition != 4 {
		t.Errorf("expected previous position 4, got %d", resultTrack.PreviousPosition)
	}

	if resultTrack.HighestPosition != 2 {
		t.Errorf("expected highest position 2, got %d", resultTrack.HighestPosition)
	}

	if resultTrack.TimesAtPeakPosition != 1 {
		t.Errorf("expected times at peak 1, got %d", resultTrack.TimesAtPeakPosition)
	}

	if resultTrack.EpisodesCount != 6 {
		t.Errorf("expected episodes count 6, got %d", resultTrack.EpisodesCount)
	}
}

func TestHandleEpisodeCreatedEvent_IncrementsTimesAtPeakWhenPeakRepeated(t *testing.T) {
	chartID := uuid.New()
	trackID := uuid.New()

	episodeRepo := &fakeEpisodeRepository{
		latestEpisode: &episode.Episode{
			ID:      uuid.New(),
			ChartID: chartID,
			Tracks: []track_episode.TrackEpisode{
				{
					TrackID:             trackID,
					CurrentPosition:     5,
					HighestPosition:     2,
					TimesAtPeakPosition: 3,
					EpisodesCount:       7,
				},
			},
		},
	}

	normalizedKey := "artist|song"

	trackRepo := &fakeEpisodeTrackRepository{
		tracksByKey: map[string]*track.Track{
			normalizedKey: {
				ID:            trackID,
				Artist:        "artist",
				Title:         "song",
				NormalizedKey: normalizedKey,
			},
		},
	}

	service := NewEpisodeService(episodeRepo, trackRepo)

	tracks := []event.TrackSnapshot{
		{Artist: "Artist", Title: "Song", CurrentPosition: 2},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 1},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    tracks,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultTrack := episodeRepo.createdEpisode.Tracks[0]

	if resultTrack.PreviousPosition != 5 {
		t.Errorf("expected previous position 5, got %d", resultTrack.PreviousPosition)
	}

	if resultTrack.HighestPosition != 2 {
		t.Errorf("expected highest position 2, got %d", resultTrack.HighestPosition)
	}

	if resultTrack.TimesAtPeakPosition != 4 {
		t.Errorf("expected times at peak 4, got %d", resultTrack.TimesAtPeakPosition)
	}

	if resultTrack.EpisodesCount != 8 {
		t.Errorf("expected episodes count 8, got %d", resultTrack.EpisodesCount)
	}
}
func TestHandleEpisodeCreatedEvent_KeepsHighestPositionWhenTrackDoesNotReachPeak(t *testing.T) {
	chartID := uuid.New()
	trackID := uuid.New()

	episodeRepo := &fakeEpisodeRepository{
		latestEpisode: &episode.Episode{
			ID:      uuid.New(),
			ChartID: chartID,
			Tracks: []track_episode.TrackEpisode{
				{
					TrackID:             trackID,
					CurrentPosition:     3,
					HighestPosition:     1,
					TimesAtPeakPosition: 2,
					EpisodesCount:       4,
				},
			},
		},
	}

	normalizedKey := "artist|song"

	trackRepo := &fakeEpisodeTrackRepository{
		tracksByKey: map[string]*track.Track{
			normalizedKey: {
				ID:            trackID,
				Artist:        "artist",
				Title:         "song",
				NormalizedKey: normalizedKey,
			},
		},
	}

	service := NewEpisodeService(episodeRepo, trackRepo)

	tracks := []event.TrackSnapshot{
		{Artist: "Artist", Title: "Song", CurrentPosition: 4},
		{Artist: "Artist 2", Title: "Track 2", CurrentPosition: 1},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 2},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 3},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    tracks,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultTrack := episodeRepo.createdEpisode.Tracks[0]

	if resultTrack.PreviousPosition != 3 {
		t.Errorf("expected previous position 3, got %d", resultTrack.PreviousPosition)
	}

	if resultTrack.HighestPosition != 1 {
		t.Errorf("expected highest position 1, got %d", resultTrack.HighestPosition)
	}

	if resultTrack.TimesAtPeakPosition != 2 {
		t.Errorf("expected times at peak 2, got %d", resultTrack.TimesAtPeakPosition)
	}

	if resultTrack.EpisodesCount != 5 {
		t.Errorf("expected episodes count 5, got %d", resultTrack.EpisodesCount)
	}
}
func TestHandleEpisodeCreatedEvent_ReturnsErrorWhenTrackDuplicatedAfterNormalization(t *testing.T) {
	chartID := uuid.New()

	episodeRepo := &fakeEpisodeRepository{}
	trackRepo := &fakeEpisodeTrackRepository{}

	service := NewEpisodeService(episodeRepo, trackRepo)

	tracks := []event.TrackSnapshot{
		{Artist: "Artist", Title: "Song", CurrentPosition: 1},
		{Artist: "artist", Title: "song", CurrentPosition: 2},
		{Artist: "Artist 3", Title: "Track 3", CurrentPosition: 3},
		{Artist: "Artist 4", Title: "Track 4", CurrentPosition: 4},
		{Artist: "Artist 5", Title: "Track 5", CurrentPosition: 5},
	}

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    tracks,
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if episodeRepo.createdEpisode != nil {
		t.Errorf("episode should not be created when duplicate track exists")
	}
}
func TestHandleEpisodeCreatedEvent_ReturnsErrorWhenGetLatestByChartIDFails(t *testing.T) {
	chartID := uuid.New()
	expectedErr := errors.New("database error")

	episodeRepo := &fakeEpisodeRepository{
		getLatestByChartIDErr: expectedErr,
	}
	trackRepo := &fakeEpisodeTrackRepository{}

	service := NewEpisodeService(episodeRepo, trackRepo)

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    validTracks(),
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if episodeRepo.createdEpisode != nil {
		t.Errorf("episode should not be created when GetLatestByChartID fails")
	}
}
func TestHandleEpisodeCreatedEvent_ReturnsErrorWhenFindOrCreateFails(t *testing.T) {
	chartID := uuid.New()
	expectedErr := errors.New("track repository error")

	episodeRepo := &fakeEpisodeRepository{}
	trackRepo := &fakeEpisodeTrackRepository{
		findOrCreateErr: expectedErr,
	}

	service := NewEpisodeService(episodeRepo, trackRepo)

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    validTracks(),
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if episodeRepo.createdEpisode != nil {
		t.Errorf("episode should not be created when FindOrCreate fails")
	}
}
func TestHandleEpisodeCreatedEvent_ReturnsErrorWhenCreateFails(t *testing.T) {
	chartID := uuid.New()
	expectedErr := errors.New("create episode failed")

	episodeRepo := &fakeEpisodeRepository{
		createErr: expectedErr,
	}

	trackRepo := &fakeEpisodeTrackRepository{}

	service := NewEpisodeService(episodeRepo, trackRepo)

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    validTracks(),
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
func TestHandleEpisodeCreatedEvent_GeneratesEpisodeAndTrackEpisodeIDs(t *testing.T) {
	chartID := uuid.New()

	episodeRepo := &fakeEpisodeRepository{}
	trackRepo := &fakeEpisodeTrackRepository{}

	service := NewEpisodeService(episodeRepo, trackRepo)

	err := service.HandleEpisodeCreatedEvent(context.Background(), event.EpisodeSnapshotEvent{
		ChartID:   chartID,
		CreatedAt: time.Now(),
		Tracks:    validTracks(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if episodeRepo.createdEpisode.ID == uuid.Nil {
		t.Errorf("expected episode ID to be generated")
	}

	for _, tr := range episodeRepo.createdEpisode.Tracks {
		if tr.ID == uuid.Nil {
			t.Errorf("expected track episode ID to be generated")
		}

		if tr.EpisodeID != episodeRepo.createdEpisode.ID {
			t.Errorf("expected track episode EpisodeID to match episode ID")
		}

		if tr.TrackID == uuid.Nil {
			t.Errorf("expected track ID to be set")
		}
	}
}
func TestGetLatestEpisodesPage_CalculatesOffset(t *testing.T) {
	repo := &fakeEpisodeRepository{
		latestPageResult: []episode.Episode{
			{ID: uuid.New(), ChartID: uuid.New()},
		},
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetLatestEpisodesPage(context.Background(), 3, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 episode, got %d", len(result))
	}

	if repo.latestPageLimit != 20 {
		t.Errorf("expected limit 20, got %d", repo.latestPageLimit)
	}

	if repo.latestPageOffset != 40 {
		t.Errorf("expected offset 40, got %d", repo.latestPageOffset)
	}
}
func TestGetLatestEpisodesPage_UsesDefaultPageWhenPageIsInvalid(t *testing.T) {
	repo := &fakeEpisodeRepository{}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	_, err := service.GetLatestEpisodesPage(context.Background(), 0, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.latestPageOffset != 0 {
		t.Errorf("expected offset 0, got %d", repo.latestPageOffset)
	}
}

func TestGetLatestEpisodesPage_UsesDefaultLimitWhenLimitIsInvalid(t *testing.T) {
	repo := &fakeEpisodeRepository{}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	_, err := service.GetLatestEpisodesPage(context.Background(), 2, 0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.latestPageLimit != 10 {
		t.Errorf("expected limit 10, got %d", repo.latestPageLimit)
	}

	if repo.latestPageOffset != 10 {
		t.Errorf("expected offset 10, got %d", repo.latestPageOffset)
	}
}
func TestGetLatestEpisodesPage_ReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeEpisodeRepository{
		latestPageErr: expectedErr,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetLatestEpisodesPage(context.Background(), 1, 10)

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

func tracksWithMeta() []dto.TrackEpisodeResponse {
	return []dto.TrackEpisodeResponse{
		{
			TrackID:             uuid.New(),
			Artist:              "artist 1",
			Title:               "song 1",
			CurrentPosition:     1,
			PreviousPosition:    2,
			HighestPosition:     1,
			TimesAtPeakPosition: 3,
			EpisodesCount:       5,
		},
		{
			TrackID:             uuid.New(),
			Artist:              "artist 2",
			Title:               "song 2",
			CurrentPosition:     2,
			PreviousPosition:    3,
			HighestPosition:     2,
			TimesAtPeakPosition: 1,
			EpisodesCount:       4,
		},
		{
			TrackID:             uuid.New(),
			Artist:              "artist 3",
			Title:               "song 3",
			CurrentPosition:     3,
			PreviousPosition:    1,
			HighestPosition:     1,
			TimesAtPeakPosition: 2,
			EpisodesCount:       6,
		},
		{
			TrackID:             uuid.New(),
			Artist:              "artist 4",
			Title:               "song 4",
			CurrentPosition:     4,
			PreviousPosition:    4,
			HighestPosition:     4,
			TimesAtPeakPosition: 1,
			EpisodesCount:       1,
		},
		{
			TrackID:             uuid.New(),
			Artist:              "artist 5",
			Title:               "song 5",
			CurrentPosition:     5,
			PreviousPosition:    0,
			HighestPosition:     5,
			TimesAtPeakPosition: 1,
			EpisodesCount:       1,
		},
	}
}
func TestGetEpisode_ReturnsEpisodeWithTracks(t *testing.T) {
	episodeID := uuid.New()
	chartID := uuid.New()
	createdAt := time.Now()

	repo := &fakeEpisodeRepository{
		episodeByID: &episode.Episode{
			ID:        episodeID,
			ChartID:   chartID,
			CreatedAt: createdAt,
		},
		tracksWithMeta: tracksWithMeta(),
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetEpisode(context.Background(), episodeID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected episode response, got nil")
	}

	if result.ID != episodeID {
		t.Errorf("expected episodeID %s, got %s", episodeID, result.ID)
	}

	if result.ChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, result.ChartID)
	}

	if !result.CreatedAt.Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, result.CreatedAt)
	}

	if len(result.Tracks) != 5 {
		t.Fatalf("expected 5 tracks, got %d", len(result.Tracks))
	}

	firstTrack := result.Tracks[0]

	if firstTrack.Artist != "Artist 1" {
		t.Errorf("expected capitalized artist %q, got %q", "Artist 1", firstTrack.Artist)
	}

	if firstTrack.Title != "Song 1" {
		t.Errorf("expected capitalized title %q, got %q", "Song 1", firstTrack.Title)
	}

	if firstTrack.CurrentPosition != 1 {
		t.Errorf("expected current position 1, got %d", firstTrack.CurrentPosition)
	}

	if firstTrack.PreviousPosition != 2 {
		t.Errorf("expected previous position 2, got %d", firstTrack.PreviousPosition)
	}

	if firstTrack.HighestPosition != 1 {
		t.Errorf("expected highest position 1, got %d", firstTrack.HighestPosition)
	}

	if firstTrack.TimesAtPeakPosition != 3 {
		t.Errorf("expected times at peak 3, got %d", firstTrack.TimesAtPeakPosition)
	}

	if firstTrack.EpisodesCount != 5 {
		t.Errorf("expected episodes count 5, got %d", firstTrack.EpisodesCount)
	}

	if firstTrack.ListenLinks.AppleMusic == "" {
		t.Errorf("expected Apple Music link")
	}

	if firstTrack.ListenLinks.YandexMusic == "" {
		t.Errorf("expected Yandex Music link")
	}
}
func TestGetEpisode_ReturnsNilWhenEpisodeNotFound(t *testing.T) {
	repo := &fakeEpisodeRepository{
		episodeByID: nil,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetEpisode(context.Background(), uuid.New())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestGetEpisode_ReturnsErrorWhenGetByIDFails(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeEpisodeRepository{
		getByIDErr: expectedErr,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetEpisode(context.Background(), uuid.New())

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
func TestGetEpisode_ReturnsErrorWhenGetTracksWithMetaFails(t *testing.T) {
	episodeID := uuid.New()
	chartID := uuid.New()

	expectedErr := errors.New("tracks loading error")

	repo := &fakeEpisodeRepository{
		episodeByID: &episode.Episode{
			ID:        episodeID,
			ChartID:   chartID,
			CreatedAt: time.Now(),
		},
		tracksWithMetaErr: expectedErr,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetEpisode(context.Background(), episodeID)

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

func TestGetEpisodesByChart_ReturnsEpisodes(t *testing.T) {
	chartID := uuid.New()

	expectedEpisodes := []episode.Episode{
		{ID: uuid.New(), ChartID: chartID},
		{ID: uuid.New(), ChartID: chartID},
	}

	repo := &fakeEpisodeRepository{
		episodesByChart: expectedEpisodes,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetEpisodesByChart(context.Background(), chartID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.requestedChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, repo.requestedChartID)
	}

	if len(result) != len(expectedEpisodes) {
		t.Fatalf("expected %d episodes, got %d", len(expectedEpisodes), len(result))
	}
}
func TestGetEpisodesByChart_ReturnsRepositoryError(t *testing.T) {
	chartID := uuid.New()
	expectedErr := errors.New("database error")

	repo := &fakeEpisodeRepository{
		getByChartIDErr: expectedErr,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetEpisodesByChart(context.Background(), chartID)

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

func TestGetLatestEpisodesWithTracks_ReturnsEpisodes(t *testing.T) {
	chartID := uuid.New()

	expectedEpisodes := []dto.EpisodeResponse{
		{
			ID:      uuid.New(),
			ChartID: chartID,
			Tracks:  tracksWithMeta(),
		},
	}

	repo := &fakeEpisodeRepository{
		latestWithTracksResult: expectedEpisodes,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetLatestEpisodesWithTracks(context.Background(), 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.latestWithTracksLimit != 10 {
		t.Errorf("expected limit 10, got %d", repo.latestWithTracksLimit)
	}

	if len(result) != len(expectedEpisodes) {
		t.Fatalf("expected %d episodes, got %d", len(expectedEpisodes), len(result))
	}
}

func TestGetLatestEpisodesWithTracks_ReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeEpisodeRepository{
		latestWithTracksErr: expectedErr,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetLatestEpisodesWithTracks(context.Background(), 10)

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
func TestGetNearestLeftEpisode_ReturnsEpisode(t *testing.T) {
	chartID := uuid.New()
	date := time.Now()

	expectedEpisode := &dto.EpisodeResponse{
		ID:        uuid.New(),
		ChartID:   chartID,
		CreatedAt: date.Add(-24 * time.Hour),
		Tracks:    tracksWithMeta(),
	}

	repo := &fakeEpisodeRepository{
		nearestLeftResult: expectedEpisode,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetNearestLeftEpisode(context.Background(), chartID, date)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expectedEpisode {
		t.Errorf("expected episode %v, got %v", expectedEpisode, result)
	}

	if repo.nearestLeftChartID != chartID {
		t.Errorf("expected chartID %s, got %s", chartID, repo.nearestLeftChartID)
	}

	if !repo.nearestLeftDate.Equal(date) {
		t.Errorf("expected date %v, got %v", date, repo.nearestLeftDate)
	}
}
func TestGetNearestLeftEpisode_ReturnsRepositoryError(t *testing.T) {
	chartID := uuid.New()
	date := time.Now()
	expectedErr := errors.New("database error")

	repo := &fakeEpisodeRepository{
		nearestLeftErr: expectedErr,
	}

	service := NewEpisodeService(repo, &fakeEpisodeTrackRepository{})

	result, err := service.GetNearestLeftEpisode(context.Background(), chartID, date)

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
