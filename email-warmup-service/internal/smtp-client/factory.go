package smtpclient

import (
	emailmanager "github.com/AbhinitKumarRai/email-warmup-service/internal/email-manager"
	usermanager "github.com/AbhinitKumarRai/email-warmup-service/internal/user-manager"
)

func CreateClient(clientType ClientType, emailManager *emailmanager.EmailManager, userManager *usermanager.UserManager) SmtpClient {
	switch clientType {
	case Google:
		return NewGoogleSmtpClient(emailManager, userManager)
	case MailPit:
		return NewMailPitSmtpClient(emailManager, userManager)
	}

	return nil
}
