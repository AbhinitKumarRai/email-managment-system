package service

import (
	"log"
	"time"

	emailmanager "github.com/AbhinitKumarRai/email-warmup-service/internal/email-manager"
	emailsender "github.com/AbhinitKumarRai/email-warmup-service/internal/email-sender"
	smtpclient "github.com/AbhinitKumarRai/email-warmup-service/internal/smtp-client"
	usermanager "github.com/AbhinitKumarRai/email-warmup-service/internal/user-manager"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
)

var UserChoosenSmtpClient smtpclient.ClientType

type EmailService struct {
	EmailManager *emailmanager.EmailManager
	UserManager  *usermanager.UserManager
	EmailSender  *emailsender.EmailSender
}

func NewEmailService(emailManager *emailmanager.EmailManager, userManager *usermanager.UserManager,
	emailSender *emailsender.EmailSender) *EmailService {
	emailService := &EmailService{
		EmailManager: emailManager,
		UserManager:  userManager,
		EmailSender:  emailSender,
	}

	return emailService
}

// SendEmail stores and sends to all users asynchronously using batch send
func (s *EmailService) SendEmail(email *model.EmailMessage) (string, error) {
	email.Timestamp = time.Now()
	s.EmailManager.AddEmail(email)
	messageId, err := s.EmailSender.SendEmailToAllUsers(email)
	if err != nil {
		log.Printf("Failed to send batch email: %v", err)
	}
	return messageId, nil
}

// AddUser adds a user and sends all stored emails to them asynchronously, rate-limited
func (s *EmailService) AddUser(user *model.User) error {
	if err := s.UserManager.AddUser(user); err != nil {
		return err
	}
	worker, err := s.EmailSender.SendEmailsToUser(user)
	if err != nil {
		return err
	}
	go func() {
		worker.Queue <- struct{}{}
	}()
	return nil
}

func (s *EmailService) GetAllUsers() ([]model.User, error) {
	return s.UserManager.GetAllUsers()
}

func (s *EmailService) GetAllEmailIds() ([]string, error) {
	return s.UserManager.GetAllEmailIds()
}

func (s *EmailService) GetAllEmails() (map[int64]model.EmailMessage, error) {
	return s.EmailManager.GetAllEmails()
}
