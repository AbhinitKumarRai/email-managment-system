package model

import "time"

type EmailEvent struct {
	MailID    string    `json:"mail_id"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
}

type EmailHealthStatus struct {
	EmailID   string    `json:"email_id"`
	Delivered int       `json:"delivered"`
	Spam      int       `json:"spam"`
	Opened    int       `json:"opened"`
	Read      int       `json:"read"`
	CheckedAt time.Time `json:"checked_at"`
}

type AggregateStats struct {
	TotalDelivered int `json:"total_delivered"`
	TotalSpam      int `json:"total_spam"`
	TotalFeedbacks int `json:"total_feedbacks"`
}
