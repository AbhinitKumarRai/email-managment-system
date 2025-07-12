package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/AbhinitKumarRai/email-health-service/internal/feedback"
	"github.com/AbhinitKumarRai/email-health-service/internal/model"
	"github.com/AbhinitKumarRai/email-health-service/internal/storage"
	sarama "github.com/IBM/sarama"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	storage *storage.Storage
}

func NewConsumer(storage *storage.Storage) *Consumer {
	return &Consumer{storage: storage}
}

func (c *Consumer) Start(brokers []string, topic string, group string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return err
	}
	ctx := context.Background()
	go func() {
		for err := range client.Errors() {
			log.Printf("Kafka error: %v", err)
		}
	}()
	return client.Consume(ctx, []string{topic}, c)
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event model.EmailEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}
		domain := strings.Split(event.To, "@")[1]
		loop, ok := feedback.Get(domain)
		if !ok {
			log.Printf("No feedback loop for domain: %s", domain)
			continue
		}
		delivered, spam, feedbacks, err := loop.Check(event.ID, event.To)
		if err != nil {
			log.Printf("Feedback check error: %v", err)
			continue
		}
		status := &model.EmailHealthStatus{
			EmailID:   event.ID,
			Delivered: 0,
			Spam:      0,
			Feedbacks: feedbacks,
			CheckedAt: event.Timestamp,
			Domain:    domain,
			Recipient: event.To,
		}
		if delivered {
			status.Delivered = 1
		}
		if spam {
			status.Spam = 1
		}
		c.storage.SaveStatus(status)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// ListenKafka starts a Kafka consumer loop (mocked for now)
func (c *Consumer) ListenKafka(brokers []string, topic string) {
	kafkaR := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "notification-service",
	})
	go func(kafkaR *kafka.Reader) {
		for {
			m, err := kafkaR.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka error: %v", err)
				continue
			}
			log.Printf("[MOCK] Received Kafka message: %s", string(m.Value))

			var notification model.Notification

			err = json.Unmarshal(m.Value, &notification)
			if err != nil {
				log.Printf("data is not correct")
			}

			if notification.Email != "" {
				s.SendEmail(notification)
			}

			if notification.Number != "" {
				s.SendMessage(notification)
			}
		}
	}(kafkaR)
}
