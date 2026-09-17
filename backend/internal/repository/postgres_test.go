package repository

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/pastebin-lite/backend/internal/model"
	"github.com/pastebin-lite/backend/internal/slug"
)

// TestPostgresRepository_Integration tests the real PostgreSQL repository against a live database.
// To run this test, provide DATABASE_URL or TEST_DATABASE_URL:
// DATABASE_URL="postgres://postgres:postgres@localhost:5432/pastebin?sslmode=disable" go test -v ./...
func TestPostgresRepository_Integration(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	if dbURL == "" {
		t.Skip("skipping PostgreSQL integration test: DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := NewPool(ctx, dbURL)
	if err != nil {
		t.Skipf("skipping PostgreSQL integration test: unable to connect to database (%v)", err)
	}
	defer pool.Close()

	// Apply migration schema to ensure table exists
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS snippets (
			id BIGSERIAL PRIMARY KEY,
			slug VARCHAR(12) UNIQUE NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	if err := RunMigration(ctx, pool, migrationSQL); err != nil {
		t.Fatalf("failed to apply migration schema: %v", err)
	}

	repo := NewPostgresRepository(pool)

	testSlug, err := slug.Generate()
	if err != nil {
		t.Fatalf("failed to generate test slug: %v", err)
	}

	snippet := &model.Snippet{
		Slug:    testSlug,
		Content: "Integration test content with real PostgreSQL",
	}

	// Step 1: Create -> Insert into PostgreSQL
	if err := repo.Create(ctx, snippet); err != nil {
		t.Fatalf("failed to create snippet: %v", err)
	}
	if snippet.ID <= 0 {
		t.Errorf("expected positive generated ID, got %d", snippet.ID)
	}
	if snippet.CreatedAt.IsZero() {
		t.Errorf("expected non-zero CreatedAt timestamp")
	}

	// Step 2: Test duplicate insert returns ErrConflict
	dupSnippet := &model.Snippet{
		Slug:    testSlug,
		Content: "duplicate slug content",
	}
	if err := repo.Create(ctx, dupSnippet); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on duplicate slug insert, got %v", err)
	}

	// Step 3: Retrieve -> Verify data
	retrieved, err := repo.GetBySlug(ctx, testSlug)
	if err != nil {
		t.Fatalf("failed to retrieve snippet: %v", err)
	}
	if retrieved.ID != snippet.ID {
		t.Errorf("expected ID %d, got %d", snippet.ID, retrieved.ID)
	}
	if retrieved.Slug != testSlug {
		t.Errorf("expected Slug %q, got %q", testSlug, retrieved.Slug)
	}
	if retrieved.Content != snippet.Content {
		t.Errorf("expected Content %q, got %q", snippet.Content, retrieved.Content)
	}

	// Step 4: Delete
	if err := repo.DeleteBySlug(ctx, testSlug); err != nil {
		t.Fatalf("failed to delete snippet: %v", err)
	}

	// Step 5: Verify deletion
	_, err = repo.GetBySlug(ctx, testSlug)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}

	// Step 6: Verify deleting already deleted snippet returns ErrNotFound
	if err := repo.DeleteBySlug(ctx, testSlug); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound deleting non-existent snippet, got %v", err)
	}
}
