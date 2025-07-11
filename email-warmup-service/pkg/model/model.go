package model

import "time"

type EmailMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

type User struct {
	Name      string    `json:"name"`
	EmailId   string    `json:"email_id"`
	CreatedAt time.Time `json:"created_at"`
}
