package smtpclient

import (
	"fmt"

	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
)

type ClientType int

const (
	MailPit ClientType = iota
	Google
)

func ParseClientType(val string) (ClientType, error) {
	switch val {
	case "MAILPIT":
		return MailPit, nil
	case "GOOGLE":
		return Google, nil
	default:
		return 0, fmt.Errorf("invalid ClientType: %s", val)
	}
}

type SmtpClient interface {
	SendEmailToMultipleReceipents(email *model.EmailMessage) (string, error)
	SendMultipleEmailsToRecipient(recipient string, rateLimitPerSecond int) error
}
