package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the response body of the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler returns a simple JSON status indicating that the service is running.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthResponse{
		Status: "ok",
	})
}
