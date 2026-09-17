package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_HealthEndpoint(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := newFakeSnippetService()
	router := NewRouter(logger, svc)

	tests := []struct {
		name           string
		method         string
		target         string
		expectedStatus int
		expectBody     bool
	}{
		{
			name:           "valid GET /health",
			method:         http.MethodGet,
			target:         "/health",
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "invalid method POST /health",
			method:         http.MethodPost,
			target:         "/health",
			expectedStatus: http.StatusMethodNotAllowed,
			expectBody:     false,
		},
		{
			name:           "not found route GET /unknown",
			method:         http.MethodGet,
			target:         "/unknown",
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:           "preflight OPTIONS request",
			method:         http.MethodOptions,
			target:         "/health",
			expectedStatus: http.StatusNoContent,
			expectBody:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, res.StatusCode)
			}

			if tc.expectBody {
				var resp HealthResponse
				if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response JSON: %v", err)
				}
				if resp.Status != "ok" {
					t.Errorf("expected status 'ok', got %q", resp.Status)
				}
			}
		})
	}
}
