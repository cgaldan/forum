package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"forum-backend/internal/domain"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{version: version}
}

// Health handles health check
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   h.version,
	}

	json.NewEncoder(w).Encode(response)
}

