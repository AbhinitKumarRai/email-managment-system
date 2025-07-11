package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AbhinitKumarRai/email-warmup-service/internal/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil || req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.RegisterEmail(req.Email); err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) BroadcastEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		From    string `json:"from"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil || req.From == "" || req.Subject == "" || req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := h.Service.BroadcastEmail(req.From, req.Subject, req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (h *Handler) ListEmails(w http.ResponseWriter, r *http.Request) {
	emails := h.Service.ListEmails()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}
