package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AbhinitKumarRai/email-health-service/internal/service"
)

type Handler struct {
	HealthService *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{HealthService: service}
}

// GET /status/{emailID}
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	emailID := r.URL.Query().Get("mail_id")

	stats, err := h.HealthService.GetStats(emailID)
	if err != nil {
		http.Error(w, "Error retrieving stats", http.StatusInternalServerError)
		return
	}

	if stats == nil {
		http.Error(w, "No stats found for this email ID", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, *stats)
}

// GET /statuses
func (h *Handler) GetAllMailStats(w http.ResponseWriter, r *http.Request) {
	allStats, err := h.HealthService.GetAllMailStats()
	if err != nil {
		http.Error(w, "Error retrieving all stats", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, allStats)
}

// Utility: JSON response writer
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
	}
}
