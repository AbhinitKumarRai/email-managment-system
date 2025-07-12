package routes

import (
	"github.com/AbhinitKumarRai/email-health-service/internal/handler"
	"github.com/AbhinitKumarRai/email-health-service/internal/service"
	"github.com/gorilla/mux"
)

func RegisterRoutes(emailService *service.Service) *mux.Router {
	router := mux.NewRouter()

	// --- Email Handlers ---

	handler := handler.NewHandler(emailService)

	router.HandleFunc("/mail_stats", handler.GetStats).Methods("GET")
	router.HandleFunc("/all_mail_stats", handler.GetAllMailStats).Methods("GET")

	return router
}
