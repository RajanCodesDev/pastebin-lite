package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/pastebin-lite/backend/internal/model"
	"github.com/pastebin-lite/backend/internal/repository"
	"github.com/pastebin-lite/backend/internal/validation"
)

// snippetService defines the operations the handler needs from the service layer.
type snippetService interface {
	CreateSnippet(ctx context.Context, input model.CreateSnippetInput) (*model.Snippet, error)
	GetSnippet(ctx context.Context, slugValue string) (*model.Snippet, error)
	DeleteSnippet(ctx context.Context, slugValue string) error
}

// SnippetHandler handles HTTP requests for snippet CRUD operations.
type SnippetHandler struct {
	svc    snippetService
	logger *slog.Logger
}

// NewSnippetHandler creates a new SnippetHandler.
func NewSnippetHandler(svc snippetService, logger *slog.Logger) *SnippetHandler {
	return &SnippetHandler{svc: svc, logger: logger}
}

// CreateSnippet handles POST /api/snippets
func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var input model.CreateSnippetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	snippet, err := h.svc.CreateSnippet(r.Context(), input)
	if err != nil {
		if errors.Is(err, validation.ErrEmptyContent) || errors.Is(err, validation.ErrContentTooLong) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.logger.Error("failed to create snippet", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create snippet")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(snippet)
}

// GetSnippet handles GET /api/snippets/{slug}
func (h *SnippetHandler) GetSnippet(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		h.respondError(w, http.StatusBadRequest, "slug is required")
		return
	}

	snippet, err := h.svc.GetSnippet(r.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.respondError(w, http.StatusNotFound, "snippet not found")
			return
		}
		h.logger.Error("failed to get snippet", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to get snippet")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snippet)
}

// DeleteSnippet handles DELETE /api/snippets/{slug}
func (h *SnippetHandler) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		h.respondError(w, http.StatusBadRequest, "slug is required")
		return
	}

	err := h.svc.DeleteSnippet(r.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.respondError(w, http.StatusNotFound, "snippet not found")
			return
		}
		h.logger.Error("failed to delete snippet", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to delete snippet")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondError writes a JSON error response with the given status code.
func (h *SnippetHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
