package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pastebin-lite/backend/internal/model"
)

func TestMemoryRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	snippet := &model.Snippet{
		Slug:      "test01",
		Content:   "fmt.Println(\"testing memory repo\")",
		CreatedAt: time.Now().UTC(),
	}

	// 1. Create
	if err := repo.Create(ctx, snippet); err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}
	if snippet.ID != 1 {
		t.Errorf("expected generated ID 1, got %d", snippet.ID)
	}

	// 2. Duplicate Create should return ErrConflict
	duplicate := &model.Snippet{
		Slug:    "test01",
		Content: "duplicate content",
	}
	if err := repo.Create(ctx, duplicate); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict on duplicate slug, got %v", err)
	}

	// 3. GetBySlug
	found, err := repo.GetBySlug(ctx, "test01")
	if err != nil {
		t.Fatalf("expected get to succeed, got %v", err)
	}
	if found.Slug != "test01" || found.Content != snippet.Content {
		t.Errorf("retrieved snippet mismatch: %+v", found)
	}

	// 4. GetBySlug not found
	_, err = repo.GetBySlug(ctx, "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nonexistent slug, got %v", err)
	}

	// 5. DeleteBySlug
	if err := repo.DeleteBySlug(ctx, "test01"); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	// 6. DeleteBySlug already deleted
	if err := repo.DeleteBySlug(ctx, "test01"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound deleting already deleted slug, got %v", err)
	}

	// 7. Verify GetBySlug returns ErrNotFound after deletion
	_, err = repo.GetBySlug(ctx, "test01")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}
}
