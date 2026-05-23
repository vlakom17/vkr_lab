package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"charts-analytics-service/internal/client/dto"

	"github.com/google/uuid"
)

type ArchiveClient struct {
	baseURL     string
	internalKey string
	client      *http.Client
}

func NewArchiveClient(baseURL, internalKey string) *ArchiveClient {
	return &ArchiveClient{
		baseURL:     baseURL,
		internalKey: internalKey,
		client:      &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ArchiveClient) GetNearestLeftEpisode(
	ctx context.Context,
	chartID uuid.UUID,
	date time.Time,
) (*dto.EpisodeResponse, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"%s/internal/episodes/nearest-left?chart_id=%s&date=%s",
			c.baseURL,
			chartID,
			date.Format(time.RFC3339),
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("archive service error: %d", resp.StatusCode)
	}

	var ep dto.EpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ep); err != nil {
		return nil, err
	}

	return &ep, nil
}

func (c *ArchiveClient) GetLatestEpisodes(
	ctx context.Context,
	limit int,
) ([]dto.EpisodeResponse, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"%s/internal/episodes/latest?limit=%d",
			c.baseURL,
			limit,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("archive service error: %d", resp.StatusCode)
	}

	var eps []dto.EpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&eps); err != nil {
		return nil, err
	}

	return eps, nil
}
