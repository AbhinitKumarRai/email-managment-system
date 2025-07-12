package main

import (
	"log"
	"net/http"
	"os"

	emailmanager "github.com/AbhinitKumarRai/email-warmup-service/internal/email-manager"
	emailsender "github.com/AbhinitKumarRai/email-warmup-service/internal/email-sender"
	"github.com/AbhinitKumarRai/email-warmup-service/internal/routes"
	"github.com/AbhinitKumarRai/email-warmup-service/internal/service"
	usermanager "github.com/AbhinitKumarRai/email-warmup-service/internal/user-manager"
)

func main() {

	usermanager := usermanager.NewUserManager()
	emailManager := emailmanager.NewEmailManager()
	emailsender := emailsender.NewEmailSender(emailManager, usermanager)
	emailService := service.NewEmailService(emailManager, usermanager, emailsender)

	// Register routes using Gorilla Mux
	router := routes.RegisterRoutes(emailService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Email warmup Service running on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
