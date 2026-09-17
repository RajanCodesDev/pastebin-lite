package service

import (
	"context"
	"errors"
	"testing"

	"github.com/pastebin-lite/backend/internal/model"
	"github.com/pastebin-lite/backend/internal/repository"
	"github.com/pastebin-lite/backend/internal/slug"
)

type fakeRepository struct {
	snippets map[string]*model.Snippet
	nextID   int64
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{snippets: make(map[string]*model.Snippet)}
}

func (f *fakeRepository) Create(_ context.Context, snippet *model.Snippet) error {
	if _, exists := f.snippets[snippet.Slug]; exists {
		return repository.ErrConflict
	}
	f.nextID++
	snippet.ID = f.nextID
	clone := *snippet
	f.snippets[snippet.Slug] = &clone
	return nil
}

func (f *fakeRepository) GetBySlug(_ context.Context, slugValue string) (*model.Snippet, error) {
	s, exists := f.snippets[slugValue]
	if !exists {
		return nil, repository.ErrNotFound
	}
	clone := *s
	return &clone, nil
}

func (f *fakeRepository) DeleteBySlug(_ context.Context, slugValue string) error {
	if _, exists := f.snippets[slugValue]; !exists {
		return repository.ErrNotFound
	}
	delete(f.snippets, slugValue)
	return nil
}

func TestCreateSnippet_Valid(t *testing.T) {
	repo := newFakeRepository()
	gen := slug.NewGenerator()
	svc := NewSnippetService(repo, gen)

	snippet, err := svc.CreateSnippet(context.Background(), model.CreateSnippetInput{Content: "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snippet.ID <= 0 {
		t.Errorf("expected positive ID, got %d", snippet.ID)
	}
	if snippet.Slug == "" {
		t.Error("expected non-empty slug")
	}
	if snippet.Content != "hello world" {
		t.Errorf("expected content 'hello world', got %q", snippet.Content)
	}
}

func TestCreateSnippet_EmptyContent(t *testing.T) {
	repo := newFakeRepository()
	svc := NewSnippetService(repo, slug.NewGenerator())

	_, err := svc.CreateSnippet(context.Background(), model.CreateSnippetInput{Content: "   "})
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}
}

func TestGetSnippet_Found(t *testing.T) {
	repo := newFakeRepository()
	svc := NewSnippetService(repo, slug.NewGenerator())

	created, err := svc.CreateSnippet(context.Background(), model.CreateSnippetInput{Content: "test content"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	snippet, err := svc.GetSnippet(context.Background(), created.Slug)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snippet.Content != "test content" {
		t.Errorf("expected content 'test content', got %q", snippet.Content)
	}
}

func TestGetSnippet_NotFound(t *testing.T) {
	repo := newFakeRepository()
	svc := NewSnippetService(repo, slug.NewGenerator())

	_, err := svc.GetSnippet(context.Background(), "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteSnippet_Found(t *testing.T) {
	repo := newFakeRepository()
	svc := NewSnippetService(repo, slug.NewGenerator())

	created, err := svc.CreateSnippet(context.Background(), model.CreateSnippetInput{Content: "delete me"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	if err := svc.DeleteSnippet(context.Background(), created.Slug); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.GetSnippet(context.Background(), created.Slug)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteSnippet_NotFound(t *testing.T) {
	repo := newFakeRepository()
	svc := NewSnippetService(repo, slug.NewGenerator())

	err := svc.DeleteSnippet(context.Background(), "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// Compile-time interface check
var _ repository.SnippetRepository = (*fakeRepository)(nil)
