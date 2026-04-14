package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AnalyticsClient struct {
	baseURL     string
	internalKey string
	httpClient  *http.Client
}

func NewAnalyticsClient(baseURL, internalKey string) *AnalyticsClient {
	return &AnalyticsClient{
		baseURL:     baseURL,
		internalKey: internalKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *AnalyticsClient) GetMostPopularChartIDs(
	ctx context.Context,
	limit int,
) ([]uuid.UUID, error) {

	url := fmt.Sprintf("%s/internal/analysis/charts/popular?limit=%d", c.baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analytics error: %d", resp.StatusCode)
	}

	var result struct {
		ChartIDs []uuid.UUID `json:"chart_ids"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.ChartIDs, nil
}

func (c *AnalyticsClient) GetUserLikedChartIDs(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {

	url := fmt.Sprintf("%s/internal/analysis/users/%s/likes", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analytics error: %d", resp.StatusCode)
	}

	var result struct {
		ChartIDs []uuid.UUID `json:"chart_ids"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.ChartIDs, nil
}

func (c *AnalyticsClient) GetUserDislikedChartIDs(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {

	url := fmt.Sprintf("%s/internal/analysis/users/%s/dislikes", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("analytics error: %d", resp.StatusCode)
	}

	var result struct {
		ChartIDs []uuid.UUID `json:"chart_ids"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.ChartIDs, nil
}
