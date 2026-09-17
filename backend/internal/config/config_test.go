package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default Port '8080', got %q", cfg.Port)
	}
	if cfg.ShutdownTimeout != 5*time.Second {
		t.Errorf("expected default ShutdownTimeout '5s', got %v", cfg.ShutdownTimeout)
	}
	if cfg.DatabaseURL != "" {
		t.Errorf("expected empty default DatabaseURL, got %q", cfg.DatabaseURL)
	}
}

func TestLoadCustom(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SHUTDOWN_TIMEOUT", "10s")
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected Port '9090', got %q", cfg.Port)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("expected ShutdownTimeout '10s', got %v", cfg.ShutdownTimeout)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("expected custom DatabaseURL, got %q", cfg.DatabaseURL)
	}
}
