package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserClient struct {
	baseURL     string
	internalKey string
	httpClient  *http.Client
}

func NewUserClient(baseURL, internalKey string) *UserClient {
	return &UserClient{
		baseURL:     baseURL,
		internalKey: internalKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *UserClient) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/internal/auth/user", c.baseURL),
		nil,
	)
	fmt.Println("URL:", fmt.Sprintf("%s/internal/auth/user", c.baseURL))
	fmt.Println("TOKEN:", token)
	if err != nil {
		return uuid.Nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return uuid.Nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return uuid.Nil, fmt.Errorf("user service error: %d", resp.StatusCode)
	}

	var result struct {
		UserID uuid.UUID `json:"user_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return uuid.Nil, err
	}
	fmt.Println("STATUS:", resp.StatusCode)
	return result.UserID, nil
}
