package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AbhinitKumarRai/email-warmup-service/internal/service"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
)

type Handler struct {
	EMailService *service.EmailService
}

func NewHandler(service *service.EmailService) *Handler {
	return &Handler{EMailService: service}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &user); err != nil || user.EmailId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.EMailService.AddUser(&user); err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var email model.EmailMessage

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &email); err != nil || email.From == "" || email.Subject == "" || email.Body == "" {
		http.Error(w, "Invalid email request", http.StatusBadRequest)
		return
	}

	messageId, err := h.EMailService.SendEmail(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare JSON response
	resp := map[string]string{
		"message_id": messageId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.EMailService.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetAllEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := h.EMailService.GetAllEmails()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}

func (h *Handler) GetAllEmailIds(w http.ResponseWriter, r *http.Request) {
	emails, err := h.EMailService.GetAllEmailIds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}
