package storage

import (
	"sync"

	"github.com/AbhinitKumarRai/email-health-service/internal/model"
)

type Storage struct {
	statuses map[int64]*model.EmailHealthStatus
	mu       sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		statuses: make(map[int64]*model.EmailHealthStatus),
	}
}

func (s *Storage) SaveStatus(status *model.EmailHealthStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.statuses[status.EmailID]; ok {
		existing.Delivered += status.Delivered
		existing.Spam += status.Spam
		existing.Feedbacks += status.Feedbacks
		existing.CheckedAt = status.CheckedAt
	} else {
		s.statuses[status.EmailID] = status
	}
}

func (s *Storage) GetStatus(emailID int64) (*model.EmailHealthStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	status, ok := s.statuses[emailID]
	return status, ok
}

func (s *Storage) GetAllStatuses() []*model.EmailHealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*model.EmailHealthStatus, 0, len(s.statuses))
	for _, v := range s.statuses {
		res = append(res, v)
	}
	return res
}

func (s *Storage) Aggregate() *model.AggregateStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	agg := &model.AggregateStats{}
	for _, v := range s.statuses {
		agg.TotalDelivered += v.Delivered
		agg.TotalSpam += v.Spam
		agg.TotalFeedbacks += v.Feedbacks
	}
	return agg
}
