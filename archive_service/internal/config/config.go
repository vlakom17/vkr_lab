package config

import (
	"os"
)

type Config struct {
	DatabaseURL    string
	KafkaBroker    string
	KafkaGroupID   string
	InternalAPIKey string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		KafkaBroker:    os.Getenv("KAFKA_BROKER"),
		KafkaGroupID:   os.Getenv("KAFKA_GROUP_ID"),
		InternalAPIKey: os.Getenv("INTERNAL_API_KEY"),
	}
}
