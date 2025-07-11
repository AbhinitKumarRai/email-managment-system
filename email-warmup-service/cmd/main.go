package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AbhinitKumarRai/email-warmup-service/internal/warmup"
)

func main() {
	service := warmup.NewService()
	handler := warmup.NewHandler(service)

	http.HandleFunc("/register", handler.RegisterEmail)
	http.HandleFunc("/broadcast", handler.BroadcastEmail)
	http.HandleFunc("/emails", handler.ListEmails)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Email Warmup Service running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
