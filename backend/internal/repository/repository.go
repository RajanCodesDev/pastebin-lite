package repository

import (
	"context"
	"errors"

	"github.com/pastebin-lite/backend/internal/model"
)

// Sentinel repository errors.
var (
	ErrNotFound = errors.New("snippet not found")
	ErrConflict = errors.New("snippet slug already exists")
)

// SnippetRepository defines the interface for snippet persistence operations.
type SnippetRepository interface {
	Create(ctx context.Context, snippet *model.Snippet) error
	GetBySlug(ctx context.Context, slug string) (*model.Snippet, error)
	DeleteBySlug(ctx context.Context, slug string) error
}
