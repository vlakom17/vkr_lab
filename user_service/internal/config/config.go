package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	RedisAddr      string
	InternalAPIKey string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		InternalAPIKey: os.Getenv("INTERNAL_API_KEY"),
	}
}
