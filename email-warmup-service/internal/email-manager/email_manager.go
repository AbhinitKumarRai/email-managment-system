package emailmanager

import (
	"sync"

	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
)

type EmailManager struct {
	emails         map[int64]model.EmailMessage
	globalRWLock   sync.RWMutex
	emailIdCounter int64
}

func NewEmailManager() *EmailManager {
	return &EmailManager{
		emails:         make(map[int64]model.EmailMessage),
		globalRWLock:   sync.RWMutex{},
		emailIdCounter: 1,
	}
}

func (e *EmailManager) AddEmail(email *model.EmailMessage) {
	e.globalRWLock.Lock()
	defer e.globalRWLock.Unlock()

	email.ID = e.emailIdCounter
	e.emailIdCounter = e.emailIdCounter + 1

	e.emails[email.ID] = *email

}

// GetAllEmails returns all stored emails as a slice
func (e *EmailManager) GetAllEmails() (map[int64]model.EmailMessage, error) {
	return e.emails, nil
}
