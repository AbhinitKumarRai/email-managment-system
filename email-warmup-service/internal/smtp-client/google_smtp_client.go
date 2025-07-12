package smtpclient

import (
	"fmt"
	"log"
	"strings"

	emailmanager "github.com/AbhinitKumarRai/email-warmup-service/internal/email-manager"
	usermanager "github.com/AbhinitKumarRai/email-warmup-service/internal/user-manager"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/smtp"
)

type GoogleSmtpClient struct {
	emailManager *emailmanager.EmailManager
	userManager  *usermanager.UserManager
}

func NewGoogleSmtpClient(emailManager *emailmanager.EmailManager, userManager *usermanager.UserManager) *GoogleSmtpClient {
	return &GoogleSmtpClient{
		emailManager: emailManager,
		userManager:  userManager,
	}
}

// SendEmailToMultipleReceipents sends a single email to multiple recipients via raw SMTP.
// Returns the generated Message-ID for tracking.
func (g *GoogleSmtpClient) SendEmailToMultipleReceipents(email *model.EmailMessage) (string, error) {
	allEmailsRegistered, err := g.userManager.GetAllEmailIds()
	if err != nil {
		return "", err
	}
	if len(allEmailsRegistered) == 0 {
		return "", nil
	}
	client, err := smtp.CreateGmailSmtpClient()
	if err != nil {
		return "", err
	}
	defer client.Quit()

	from := email.From
	if err := client.Mail(from); err != nil {
		return "", err
	}
	for _, recipient := range allEmailsRegistered {
		if err := client.Rcpt(recipient); err != nil {
			log.Printf("[SMTP] Error adding recipient %s: %v", recipient, err)
			continue
		}
	}

	wc, err := client.Data()
	if err != nil {
		return "", err
	}

	// Generate a proper RFC 5322-compliant Message-ID
	messageID := fmt.Sprintf("email-%d@gmail.com", email.ID)

	toHeader := strings.Join(allEmailsRegistered, ", ")
	msg := fmt.Sprintf(
		"Subject: %s\r\nTo: %s\r\nFrom: %s\r\nDate: %s\r\nMessage-ID: %s\r\n\r\n%s\r\n",
		email.Subject,
		toHeader,
		email.From,
		email.Timestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		messageID,
		email.Body,
	)

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	log.Printf("✅ Sent email ID %d with Message-ID %s to %d users", email.ID, messageID, len(allEmailsRegistered))
	return messageID, nil
}

func (g *GoogleSmtpClient) SendMultipleEmailsToRecipient(recipient string, rateLimitPerSecond int) error {

	emails, err := g.emailManager.GetAllEmails()
	if err != nil {
		return err
	}

	if len(emails) == 0 {
		return nil
	}

	client, err := smtp.CreateGmailSmtpClient()
	if err != nil {
		return err
	}
	defer client.Quit()

	sentCount := 0

	for _, email := range emails {
		from := email.From
		if err := client.Mail(from); err != nil {
			log.Printf("[SMTP] Error MAIL FROM for %s: %v", from, err)
			continue
		}
		if err := client.Rcpt(recipient); err != nil {
			log.Printf("[SMTP] Error RCPT TO for %s: %v", recipient, err)
			continue
		}
		wc, err := client.Data()
		if err != nil {
			log.Printf("[SMTP] Error DATA for %s: %v", recipient, err)
			continue
		}
		msg := fmt.Sprintf("Subject: %s\r\nTo: %s\r\nFrom: %s\r\nDate: %s\r\nMessage-ID: <%d@localhost>\r\n\r\n%s\r\n",
			email.Subject, recipient, email.From, email.Timestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700"), email.ID, email.Body)
		_, err = wc.Write([]byte(msg))
		if err != nil {
			log.Printf("[SMTP] Error writing DATA for %s: %v", recipient, err)
			wc.Close()
			client.Reset()
			continue
		}
		if err := wc.Close(); err != nil {
			log.Printf("[SMTP] Error closing DATA for %s: %v", recipient, err)
			client.Reset()
			continue
		}
		sentCount++

		if err := client.Reset(); err != nil {
			log.Printf("[SMTP] Error resetting client: %v", err)
			continue
		}
	}
	if sentCount > 0 {
		log.Printf("✅ Sent %d emails to %s", sentCount, recipient)
	}
	return nil
}
