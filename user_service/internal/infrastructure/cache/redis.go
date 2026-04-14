package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisAddr string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
