package config

import (
	"os"
)

type Config struct {
	DatabaseURL         string
	KafkaBroker         string
	UserServiceURL      string
	AnalyticsServiceURL string
	InternalAPIKey      string
}

func Load() *Config {
	return &Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		KafkaBroker:         os.Getenv("KAFKA_BROKER"),
		UserServiceURL:      os.Getenv("USER_SERVICE_URL"),
		AnalyticsServiceURL: os.Getenv("ANALYTICS_SERVICE_URL"),
		InternalAPIKey:      os.Getenv("INTERNAL_API_KEY"),
	}
}
