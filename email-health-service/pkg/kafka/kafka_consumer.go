package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/AbhinitKumarRai/email-health-service/pkg/model"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	eventChan chan model.EmailEvent
}

func NewConsumer(eventChan chan model.EmailEvent) *Consumer {
	c := &Consumer{
		eventChan: eventChan,
	}
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	go c.ListenEmailEvents(brokers, topic)
	return c
}

func (c *Consumer) ListenEmailEvents(brokers []string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "health-service-email",
	})
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}
		var event model.EmailEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		c.eventChan <- event // Push to channel
	}
}
