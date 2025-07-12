package feedback

import "github.com/AbhinitKumarRai/email-health-service/pkg/model"

type DomainType int

const (
	Google DomainType = iota
	Yahoo
)

type FeedbackLoop interface {
	RegisterMailID(emailID string)
	GetAllStats() map[string]*model.EmailHealthStatus
	GetStats(mailId string) *model.EmailHealthStatus
}
