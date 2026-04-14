package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	client *redis.Client
}

func NewAuthRepository(client *redis.Client) *AuthRepository {
	return &AuthRepository{
		client: client,
	}
}

func (r *AuthRepository) Create(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error {

	err := r.client.Set(ctx, token, userID.String(), ttl).Err()
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "user:"+userID.String(), token, ttl).Err()
}

func (r *AuthRepository) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	val, err := r.client.Get(ctx, token).Result()
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(val)
}

func (r *AuthRepository) Delete(ctx context.Context, token string) error {

	userID, err := r.GetUserID(ctx, token)
	if err != nil {
		return err
	}

	r.client.Del(ctx, "user:"+userID.String())

	return r.client.Del(ctx, token).Err()
}

func (r *AuthRepository) GetTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error) {

	key := "user:" + userID.String()

	token, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	return token, err
}
