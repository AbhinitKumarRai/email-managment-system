package model

import "time"

type EmailEvent struct {
	ID        int64     `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

type EmailHealthStatus struct {
	EmailID   int64     `json:"email_id"`
	Delivered int       `json:"delivered"`
	Spam      int       `json:"spam"`
	Feedbacks int       `json:"feedbacks"`
	CheckedAt time.Time `json:"checked_at"`
	Domain    string    `json:"domain"`
	Recipient string    `json:"recipient"`
}

type AggregateStats struct {
	TotalDelivered int `json:"total_delivered"`
	TotalSpam      int `json:"total_spam"`
	TotalFeedbacks int `json:"total_feedbacks"`
}
