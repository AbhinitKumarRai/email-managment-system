package smtpclient

import (
	"fmt"
	"log"
	"strings"
	"time"

	emailmanager "github.com/AbhinitKumarRai/email-warmup-service/internal/email-manager"
	usermanager "github.com/AbhinitKumarRai/email-warmup-service/internal/user-manager"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/smtp"
)

type MailPitClient struct {
	emailManager *emailmanager.EmailManager
	userManager  *usermanager.UserManager
}

func NewMailPitSmtpClient(emailManager *emailmanager.EmailManager, userManager *usermanager.UserManager) *GoogleSmtpClient {
	return &GoogleSmtpClient{
		emailManager: emailManager,
		userManager:  userManager,
	}
}

// SendEmailToMultipleReceipents sends a single email to multiple recipients via raw SMTP (MailPit).
// It returns the generated Message-ID for tracking.
func (m *MailPitClient) SendEmailToMultipleReceipents(email *model.EmailMessage) (string, error) {
	allEmailsRegistered, err := m.userManager.GetAllEmailIds()
	if err != nil {
		return "", err
	}
	if len(allEmailsRegistered) == 0 {
		return "", nil
	}
	conn, err := smtp.CreateSmtpClient()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	readResponse := func() {
		buf := make([]byte, 1024)
		_, _ = conn.Read(buf) // ignore output for now
	}
	write := func(format string, args ...interface{}) {
		fmt.Fprintf(conn, format, args...)
	}

	readResponse() // banner
	write("HELO localhost\r\n")
	readResponse()
	write("MAIL FROM:<%s>\r\n", email.From)
	readResponse()

	for _, recipient := range allEmailsRegistered {
		write("RCPT TO:<%s>\r\n", recipient)
		readResponse()
	}
	write("DATA\r\n")
	readResponse()

	messageID := fmt.Sprintf("email-%d@mailpit.com", email.ID)

	toHeader := strings.Join(allEmailsRegistered, ", ")
	headers := fmt.Sprintf(
		"Subject: %s\r\nTo: %s\r\nFrom: %s\r\nDate: %s\r\nMessage-ID: %s\r\n\r\n",
		email.Subject,
		toHeader,
		email.From,
		email.Timestamp.Format(time.RFC1123Z),
		messageID,
	)
	write(headers + email.Body + "\r\n.\r\n")
	readResponse()

	write("QUIT\r\n")
	readResponse()

	log.Printf("✅ Sent email %d (Message-ID: %s) to %d recipients", email.ID, messageID, len(allEmailsRegistered))
	return messageID, nil
}

func (m *MailPitClient) SendMultipleEmailsToRecipient(recipient string, rateLimitPerSecond int) error {

	emails, err := m.emailManager.GetAllEmails()
	if err != nil {
		return err
	}

	if len(emails) == 0 {
		return nil
	}
	conn, err := smtp.CreateSmtpClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	readResponse := func() {
		buf := make([]byte, 1024)
		conn.Read(buf) // ignore output
	}
	write := func(format string, args ...interface{}) {
		fmt.Fprintf(conn, format, args...)
	}

	readResponse() // banner
	write("HELO localhost\r\n")
	readResponse()

	sentCount := 0
	for _, email := range emails {
		write("MAIL FROM:<%s>\r\n", email.From)
		readResponse()
		write("RCPT TO:<%s>\r\n", recipient)
		readResponse()
		write("DATA\r\n")
		readResponse()
		body := fmt.Sprintf("Subject: %s\r\nTo: %s\r\nFrom: %s\r\nDate: %s\r\nMessage-ID: <%d@localhost>\r\n\r\n%s\r\n.\r\n",
			email.Subject, recipient, email.From, email.Timestamp.Format(time.RFC1123Z), email.ID, email.Body)
		write(body)
		readResponse()
		sentCount++
	}
	write("QUIT\r\n")
	readResponse()
	if sentCount > 0 {
		log.Printf("✅ Sent %d emails to %s", sentCount, recipient)
	}
	return nil
}
