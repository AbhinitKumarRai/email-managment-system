package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/AbhinitKumarRai/email-health-service/internal/feedback"
	"github.com/AbhinitKumarRai/email-health-service/internal/storage"
)

func main() {
	store := storage.NewStorage()
	// Optionally start Kafka consumer if env set
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")
	if brokers != "" && topic != "" {
		log.Printf("Starting Kafka consumer for topic %s", topic)
		brokerList := strings.Split(brokers, ",")
		svc.ListenKafka(brokerList, topic)
	}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		emailIDStr := r.URL.Query().Get("id")
		emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		status, ok := store.GetStatus(emailID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	http.HandleFunc("/statuses", func(w http.ResponseWriter, r *http.Request) {
		statuses := store.GetAllStatuses()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	})

	http.HandleFunc("/aggregate", func(w http.ResponseWriter, r *http.Request) {
		agg := store.Aggregate()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(agg)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Email Health Service running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
