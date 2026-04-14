package config

import (
	"os"
)

type Config struct {
	DatabaseURL       string
	KafkaBroker       string
	KafkaGroupID      string
	UserServiceURL    string
	ArchiveServiceURL string
	InternalAPIKey    string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		KafkaBroker:       os.Getenv("KAFKA_BROKER"),
		KafkaGroupID:      os.Getenv("KAFKA_GROUP_ID"),
		UserServiceURL:    os.Getenv("USER_SERVICE_URL"),
		ArchiveServiceURL: os.Getenv("ARCHIVE_SERVICE_URL"),
		InternalAPIKey:    os.Getenv("INTERNAL_API_KEY"),
	}
}
