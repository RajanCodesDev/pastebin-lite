package handler

import (
	"log/slog"
	"net/http"
	"time"
)

// statusWriter wraps http.ResponseWriter to capture the HTTP status code.
type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs HTTP request details including method, path, status, and duration.
func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(sw, r)

		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.statusCode,
			"duration", time.Since(start).String(),
		)
	})
}

// recoveryMiddleware catches panics and returns a 500 Internal Server Error.
func recoveryMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered", "error", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adds basic CORS headers to allow cross-origin requests from the frontend.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewRouter registers routes and returns the configured root http.Handler.
func NewRouter(logger *slog.Logger, snippetSvc snippetService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", HealthHandler)

	snippetHandler := NewSnippetHandler(snippetSvc, logger)
	mux.HandleFunc("POST /api/snippets", snippetHandler.CreateSnippet)
	mux.HandleFunc("GET /api/snippets/{slug}", snippetHandler.GetSnippet)
	mux.HandleFunc("DELETE /api/snippets/{slug}", snippetHandler.DeleteSnippet)

	var handler http.Handler = mux
	handler = corsMiddleware(handler)
	handler = loggingMiddleware(logger, handler)
	handler = recoveryMiddleware(logger, handler)

	return handler
}
