package emailsender

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	emailmanager "github.com/AbhinitKumarRai/email-warmup-service/internal/email-manager"
	smtpclient "github.com/AbhinitKumarRai/email-warmup-service/internal/smtp-client"
	usermanager "github.com/AbhinitKumarRai/email-warmup-service/internal/user-manager"
	kafkaPkg "github.com/AbhinitKumarRai/email-warmup-service/pkg/kafka"
	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
	"github.com/segmentio/kafka-go"
)

var UserChoosenSmtpClient smtpclient.ClientType

type UserWorker struct {
	ratePerMinute int
	Queue         chan struct{}
	Stop          chan struct{}
}

type EmailSender struct {
	EmailManager       *emailmanager.EmailManager
	UserManager        *usermanager.UserManager
	SmtpClientRegister map[smtpclient.ClientType]smtpclient.SmtpClient
	KafkaWriteInst     map[string]*kafka.Writer
	domainLimits       map[string]int
	userWorkers        map[string]*UserWorker
	workerMutex        sync.Mutex
}

func NewEmailSender(emailManager *emailmanager.EmailManager, userManager *usermanager.UserManager) *EmailSender {
	emailSender := &EmailSender{
		EmailManager:       emailManager,
		UserManager:        userManager,
		domainLimits:       defaultDomainLimits(),
		userWorkers:        make(map[string]*UserWorker),
		SmtpClientRegister: make(map[smtpclient.ClientType]smtpclient.SmtpClient),
		KafkaWriteInst:     make(map[string]*kafka.Writer),
	}

	emailEventwriter, err := kafkaPkg.ConnectToKafkaWriterForTopic(kafkaPkg.EmailEventTopic)
	if err != nil {
		panic("unable to create kafka writer")
	}

	emailSender.KafkaWriteInst[kafkaPkg.EmailEventTopic] = emailEventwriter
	if err := emailSender.InitSmtpClients(); err != nil {
		log.Fatalf("Invalid CHOOSEN_SMTP_CLIENT_TYPE: %v", err)
	}

	return emailSender
}

func (s *EmailSender) InitSmtpClients() error {
	s.SmtpClientRegister = map[smtpclient.ClientType]smtpclient.SmtpClient{
		smtpclient.Google:  smtpclient.CreateClient(smtpclient.Google, s.EmailManager, s.UserManager),
		smtpclient.MailPit: smtpclient.CreateClient(smtpclient.MailPit, s.EmailManager, s.UserManager),
	}

	clientTypeStr := os.Getenv("CHOOSEN_SMTP_CLIENT_TYPE")
	clientType, err := smtpclient.ParseClientType(clientTypeStr)
	if err != nil {
		return err
	}

	UserChoosenSmtpClient = clientType
	return nil
}

func defaultDomainLimits() map[string]int {
	return map[string]int{
		"gmail.com":   10, // 10 emails/minute
		"yahoo.com":   5,
		"outlook.com": 8,
		"default":     5,
	}
}

func (s *EmailSender) getDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "default"
	}
	domain := strings.ToLower(parts[1])
	if _, ok := s.domainLimits[domain]; ok {
		return domain
	}
	return "default"
}

func (s *EmailSender) SendEmailToAllUsers(email *model.EmailMessage) (string, error) {
	messageId, err := s.SmtpClientRegister[UserChoosenSmtpClient].SendEmailToMultipleReceipents(email)
	if err != nil {
		return "", err
	}

	event := map[string]interface{}{
		"mail_id":    messageId,
		"subject":    email.Subject,
		"created_at": email.Timestamp,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("failed to marshal kafka event: %w", err)
	}
	err = kafkaPkg.SendMessageToTopic(kafkaPkg.EmailEventTopic, data, s.KafkaWriteInst[kafkaPkg.EmailEventTopic])
	if err != nil {
		return "", err
	}

	return messageId, nil
}

func (s *EmailSender) SendEmailsToUser(user *model.User) (*UserWorker, error) {
	s.workerMutex.Lock()
	defer s.workerMutex.Unlock()

	w, ok := s.userWorkers[user.EmailId]
	if !ok {
		domain := s.getDomain(user.EmailId)
		rate := s.domainLimits[domain]
		w = &UserWorker{
			ratePerMinute: rate,
			Queue:         make(chan struct{}, 10), // just a signal channel
			Stop:          make(chan struct{}),
		}
		s.userWorkers[user.EmailId] = w
		go s.runWorker(user, w, UserChoosenSmtpClient)
	}
	return w, nil
}

func (s *EmailSender) runWorker(user *model.User, w *UserWorker, smtpClient smtpclient.ClientType) {
	for {
		select {
		case <-w.Queue:
			domain := s.getDomain(user.EmailId)
			rate := s.domainLimits[domain]
			if rate <= 0 {
				rate = 1
			}
			err := s.SmtpClientRegister[smtpClient].SendMultipleEmailsToRecipient(user.EmailId, rate)
			if err != nil {
				log.Printf("Failed to send batch emails to %s: %v", user.EmailId, err)
			}
		case <-w.Stop:
			return
		}
	}
}
