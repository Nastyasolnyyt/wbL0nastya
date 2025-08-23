package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	DBURL        string
	KafkaBrokers []string
	KafkaTopic   string
	HTTPPort     string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, errors.New("DB_URL environment variable is not set")
	}

	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersStr == "" {
		return nil, errors.New("KAFKA_BROKERS environment variable is not set")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		return nil, errors.New("KAFKA_TOPIC environment variable is not set")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // Default port
	}

	return &Config{
		DBURL:        dbURL,
		KafkaBrokers: strings.Split(kafkaBrokersStr, ","),
		KafkaTopic:   kafkaTopic,
		HTTPPort:     httpPort,
	}, nil
}
