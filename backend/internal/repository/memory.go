package repository

import (
	"context"
	"sync"
	"time"

	"github.com/pastebin-lite/backend/internal/model"
)

// MemoryRepository provides a thread-safe in-memory implementation of SnippetRepository,
// useful for unit tests and local mock development without requiring a live PostgreSQL instance.
type MemoryRepository struct {
	mu       sync.RWMutex
	snippets map[string]*model.Snippet
	nextID   int64
}

// NewMemoryRepository initializes an empty in-memory snippet repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		snippets: make(map[string]*model.Snippet),
		nextID:   1,
	}
}

// Create inserts a snippet or returns ErrConflict if the slug already exists.
func (r *MemoryRepository) Create(_ context.Context, snippet *model.Snippet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.snippets[snippet.Slug]; exists {
		return ErrConflict
	}

	snippet.ID = r.nextID
	r.nextID++
	if snippet.CreatedAt.IsZero() {
		snippet.CreatedAt = time.Now().UTC()
	}

	// Store copy
	clone := *snippet
	r.snippets[snippet.Slug] = &clone

	return nil
}

// GetBySlug retrieves a snippet by slug or returns ErrNotFound.
func (r *MemoryRepository) GetBySlug(_ context.Context, slug string) (*model.Snippet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snippet, exists := r.snippets[slug]
	if !exists {
		return nil, ErrNotFound
	}

	clone := *snippet
	return &clone, nil
}

// DeleteBySlug removes a snippet by slug or returns ErrNotFound.
func (r *MemoryRepository) DeleteBySlug(_ context.Context, slug string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.snippets[slug]; !exists {
		return ErrNotFound
	}

	delete(r.snippets, slug)
	return nil
}
