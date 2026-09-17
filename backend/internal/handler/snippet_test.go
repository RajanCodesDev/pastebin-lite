package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pastebin-lite/backend/internal/model"
	"github.com/pastebin-lite/backend/internal/repository"
	"github.com/pastebin-lite/backend/internal/slug"
	"github.com/pastebin-lite/backend/internal/validation"
)

// fakeSnippetService implements the snippetService interface for handler tests.
type fakeSnippetService struct {
	snippets map[string]*model.Snippet
	nextID   int64
}

func newFakeSnippetService() *fakeSnippetService {
	return &fakeSnippetService{snippets: make(map[string]*model.Snippet)}
}

func (f *fakeSnippetService) CreateSnippet(_ context.Context, input model.CreateSnippetInput) (*model.Snippet, error) {
	if err := validation.ValidateCreateSnippetInput(input); err != nil {
		return nil, err
	}
	s, err := slug.Generate()
	if err != nil {
		return nil, err
	}
	for _, sn := range f.snippets {
		if sn.Slug == s {
			return f.CreateSnippet(context.Background(), input)
		}
	}
	f.nextID++
	snippet := &model.Snippet{ID: f.nextID, Slug: s, Content: input.Content}
	clone := *snippet
	f.snippets[s] = &clone
	return &clone, nil
}

func (f *fakeSnippetService) GetSnippet(_ context.Context, slugValue string) (*model.Snippet, error) {
	s, exists := f.snippets[slugValue]
	if !exists {
		return nil, repository.ErrNotFound
	}
	clone := *s
	return &clone, nil
}

func (f *fakeSnippetService) DeleteSnippet(_ context.Context, slugValue string) error {
	if _, exists := f.snippets[slugValue]; !exists {
		return repository.ErrNotFound
	}
	delete(f.snippets, slugValue)
	return nil
}

func TestCreateSnippetHandler_Success(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	body := `{"content":"hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/api/snippets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	var resp model.Snippet
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Content != "hello world" {
		t.Errorf("expected content 'hello world', got %q", resp.Content)
	}
	if resp.Slug == "" {
		t.Error("expected non-empty slug")
	}
}

func TestCreateSnippetHandler_EmptyContent(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	body := `{"content":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/snippets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestGetSnippetHandler_Success(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	created, err := svc.CreateSnippet(context.Background(), model.CreateSnippetInput{Content: "get test"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/snippets/"+created.Slug, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var resp model.Snippet
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Content != "get test" {
		t.Errorf("expected content 'get test', got %q", resp.Content)
	}
}

func TestGetSnippetHandler_NotFound(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/snippets/nonexistent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestDeleteSnippetHandler_Success(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	created, err := svc.CreateSnippet(context.Background(), model.CreateSnippetInput{Content: "delete me"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/snippets/"+created.Slug, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", res.StatusCode)
	}
}

func TestDeleteSnippetHandler_NotFound(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/snippets/nonexistent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestSnippetHandler_InvalidJSON(t *testing.T) {
	svc := newFakeSnippetService()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := NewRouter(logger, svc)

	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/snippets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// Verify fake implements the interface
var _ snippetService = (*fakeSnippetService)(nil)
