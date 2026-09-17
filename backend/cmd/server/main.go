package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pastebin-lite/backend/internal/config"
	"github.com/pastebin-lite/backend/internal/handler"
	"github.com/pastebin-lite/backend/internal/repository"
	"github.com/pastebin-lite/backend/internal/service"
	"github.com/pastebin-lite/backend/internal/slug"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Initialize slug generator
	slugGen := slug.NewGenerator()

	// Initialize repository and service
	var snippetSvc *service.SnippetService
	if cfg.DatabaseURL != "" {
		pool, err := initPostgresPool(logger, cfg.DatabaseURL)
		if err != nil {
			logger.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		repo := repository.NewPostgresRepository(pool)
		snippetSvc = service.NewSnippetService(repo, slugGen)
	} else {
		// In-memory repository for development/testing without PostgreSQL
		logger.Warn("DATABASE_URL not set, using in-memory repository")
		repo := repository.NewMemoryRepository()
		snippetSvc = service.NewSnippetService(repo, slugGen)
	}

	// Initialize router and routes
	router := handler.NewRouter(logger, snippetSvc)

	// Configure HTTP server with production-grade timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to catch server startup errors
	serverErrors := make(chan error, 1)

	// Start server in background goroutine
	go func() {
		logger.Info("starting HTTP server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Listen for OS interrupt and termination signals for graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Block until a shutdown signal is received or the server encounters an error
	select {
	case err := <-serverErrors:
		logger.Error("server startup failed", "error", err)
		os.Exit(1)
	case <-shutdownCtx.Done():
		logger.Info("shutdown signal received, initiating graceful shutdown")

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server forced to shutdown", "error", err)
			_ = server.Close()
			os.Exit(1)
		}

		logger.Info("server exited cleanly")
	}
}

func initPostgresPool(logger *slog.Logger, databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := repository.NewPool(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database pool: %w", err)
	}

	// Run migration to ensure table exists
	migrationSQL, err := readMigrationFile()
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to read migration file: %w", err)
	}

	if err := repository.RunMigration(ctx, pool, string(migrationSQL)); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	logger.Info("database connection established and migration applied")
	return pool, nil
}

func readMigrationFile() (string, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(currentFile)
	// Navigate from cmd/server/ to backend/migrations/
	migrationPath := filepath.Join(baseDir, "..", "..", "migrations", "001_create_snippets.up.sql")
	data, err := os.ReadFile(migrationPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
