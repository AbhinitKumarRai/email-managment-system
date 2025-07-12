package model

import "time"

type EmailMessage struct {
	ID        int64     `json:"id"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

type User struct {
	Name      string    `json:"name"`
	EmailId   string    `json:"email_id"`
	CreatedAt time.Time `json:"created_at"`
}
