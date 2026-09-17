package service

import (
	"context"
	"fmt"

	"github.com/pastebin-lite/backend/internal/model"
	"github.com/pastebin-lite/backend/internal/repository"
	"github.com/pastebin-lite/backend/internal/slug"
	"github.com/pastebin-lite/backend/internal/validation"
)

// SnippetService orchestrates snippet operations across validation, slug generation, and persistence.
type SnippetService struct {
	repo    repository.SnippetRepository
	slugGen slug.Generator
}

// NewSnippetService creates a SnippetService with the given repository and slug generator.
func NewSnippetService(repo repository.SnippetRepository, slugGen slug.Generator) *SnippetService {
	return &SnippetService{repo: repo, slugGen: slugGen}
}

// CreateSnippet validates input, generates a slug, persists the snippet, and returns it.
func (s *SnippetService) CreateSnippet(ctx context.Context, input model.CreateSnippetInput) (*model.Snippet, error) {
	if err := validation.ValidateCreateSnippetInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	slugValue, err := s.slugGen.Generate()
	if err != nil {
		return nil, fmt.Errorf("slug generation failed: %w", err)
	}

	snippet := &model.Snippet{
		Slug:    slugValue,
		Content: input.Content,
	}

	if err := s.repo.Create(ctx, snippet); err != nil {
		return nil, fmt.Errorf("failed to create snippet: %w", err)
	}

	return snippet, nil
}

// GetSnippet retrieves a snippet by its slug.
func (s *SnippetService) GetSnippet(ctx context.Context, slugValue string) (*model.Snippet, error) {
	snippet, err := s.repo.GetBySlug(ctx, slugValue)
	if err != nil {
		return nil, err
	}
	return snippet, nil
}

// DeleteSnippet removes a snippet by its slug.
func (s *SnippetService) DeleteSnippet(ctx context.Context, slugValue string) error {
	if err := s.repo.DeleteBySlug(ctx, slugValue); err != nil {
		return err
	}
	return nil
}
