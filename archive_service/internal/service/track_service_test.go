package service

import (
	"context"
	"errors"
	"testing"

	"charts-archive-service/internal/domain/track"
)

type fakeTrackRepository struct {
	searchQuery  string
	searchCalled bool
	searchResult []track.Track
	searchErr    error
}

func (r *fakeTrackRepository) FindOrCreate(ctx context.Context, artist, title, normalizedKey string) (*track.Track, error) {
	return nil, nil
}

func (r *fakeTrackRepository) Search(ctx context.Context, query string) ([]track.Track, error) {
	r.searchCalled = true
	r.searchQuery = query

	if r.searchErr != nil {
		return nil, r.searchErr
	}

	return r.searchResult, nil
}

func TestSearchTracks_ReturnsEmptyWhenQueryIsTooShort(t *testing.T) {
	repo := &fakeTrackRepository{}

	service := NewTrackService(repo)

	result, err := service.SearchTracks(context.Background(), "a")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d tracks", len(result))
	}

	if repo.searchCalled {
		t.Errorf("Search should not be called when query is too short")
	}
}

func TestSearchTracks_TrimsQueryAndSearches(t *testing.T) {
	expectedTracks := []track.Track{
		{
			Artist: "Artist 1",
			Title:  "Track 1",
		},
	}

	repo := &fakeTrackRepository{
		searchResult: expectedTracks,
	}

	service := NewTrackService(repo)

	result, err := service.SearchTracks(context.Background(), "  track  ")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.searchCalled {
		t.Fatalf("expected Search to be called")
	}

	if repo.searchQuery != "track" {
		t.Errorf("expected query %q, got %q", "track", repo.searchQuery)
	}

	if len(result) != len(expectedTracks) {
		t.Fatalf("expected %d tracks, got %d", len(expectedTracks), len(result))
	}
}

func TestSearchTracks_ReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeTrackRepository{
		searchErr: expectedErr,
	}

	service := NewTrackService(repo)

	result, err := service.SearchTracks(context.Background(), "track")

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
