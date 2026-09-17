package config

import (
	"os"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Port            string
	ShutdownTimeout time.Duration
	DatabaseURL     string
}

// Load loads configuration from environment variables with sensible defaults.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	shutdownTimeoutStr := os.Getenv("SHUTDOWN_TIMEOUT")
	shutdownTimeout := 5 * time.Second
	if shutdownTimeoutStr != "" {
		if parsed, err := time.ParseDuration(shutdownTimeoutStr); err == nil {
			shutdownTimeout = parsed
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")

	return &Config{
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
		DatabaseURL:     databaseURL,
	}
}
