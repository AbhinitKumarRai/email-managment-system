package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/AbhinitKumarRai/email-health-service/internal/feedback"
	"github.com/AbhinitKumarRai/email-health-service/internal/routes"
	"github.com/AbhinitKumarRai/email-health-service/internal/service"
	kafkaPkg "github.com/AbhinitKumarRai/email-health-service/pkg/kafka"
	"github.com/AbhinitKumarRai/email-health-service/pkg/model"
)

func main() {

	var emailEventChan = make(chan model.EmailEvent, 100)

	service := service.NewService(emailEventChan)

	kafkaPkg.NewConsumer(emailEventChan)

	// Register routes using Gorilla Mux
	router := routes.RegisterRoutes(service)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Email Health Service running on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
