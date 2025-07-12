package routes

import (
	"github.com/AbhinitKumarRai/email-warmup-service/internal/handler"
	"github.com/AbhinitKumarRai/email-warmup-service/internal/service"
	"github.com/gorilla/mux"
)

func RegisterRoutes(emailService *service.EmailService) *mux.Router {
	router := mux.NewRouter()

	// --- Email Handlers ---

	handler := handler.NewHandler(emailService)

	router.HandleFunc("/register", handler.RegisterUser).Methods("POST")
	router.HandleFunc("/send", handler.SendEmail).Methods("POST")
	router.HandleFunc("/users", handler.ListUsers).Methods("GET")
	router.HandleFunc("/emails", handler.GetAllEmails).Methods("GET")
	router.HandleFunc("/emailIds", handler.GetAllEmailIds).Methods("GET")

	return router
}
