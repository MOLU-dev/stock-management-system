package config


import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	Environment   string
	JWTSecret     string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgresql://molu:incorrect@localhost:5432/sms?sslmode=disable"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
