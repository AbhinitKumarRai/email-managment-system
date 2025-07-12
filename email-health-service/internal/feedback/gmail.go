package feedback

import (
	"math/rand"
	"sync"
	"time"

	"github.com/AbhinitKumarRai/email-health-service/pkg/model"
)

type GmailFeedbackLoop struct {
	stats  map[string]*model.EmailHealthStatus
	ticker *time.Ticker
	mu     sync.RWMutex
}

func NewGmailFeedbackLoop() *GmailFeedbackLoop {
	gmailFeedbackInst := &GmailFeedbackLoop{
		stats:  make(map[string]*model.EmailHealthStatus),
		ticker: time.NewTicker(5 * time.Second), // simulate every 5s
	}

	go gmailFeedbackInst.simulateFeedback()
	return gmailFeedbackInst
}

func (g *GmailFeedbackLoop) RegisterMailID(emailID string) {
	g.mu.Lock()
	g.stats[emailID] = &model.EmailHealthStatus{
		EmailID: emailID,
	}

	g.mu.Unlock()
}

func (g *GmailFeedbackLoop) simulateFeedback() {
	for {
		select {
		case <-g.ticker.C:
			g.mu.Lock()
			for _, stat := range g.stats {
				stat.Delivered += rand.Intn(100) + 50
				stat.Spam += rand.Intn(3)
				stat.Opened += rand.Intn(80) + 20
				stat.Read += rand.Intn(60)
				stat.CheckedAt = time.Now()
			}
			g.mu.Unlock()
		}
	}
}

// Expose current snapshot of statuses (deep copy)
func (g *GmailFeedbackLoop) GetAllStats() map[string]*model.EmailHealthStatus {
	g.mu.RLock()
	defer g.mu.RUnlock()

	copyMap := make(map[string]*model.EmailHealthStatus)
	for id, s := range g.stats {
		copyMap[id] = &model.EmailHealthStatus{
			EmailID:   s.EmailID,
			Delivered: s.Delivered,
			Spam:      s.Spam,
			Opened:    s.Opened,
			Read:      s.Read,
			CheckedAt: s.CheckedAt,
		}
	}
	return copyMap
}

func (g *GmailFeedbackLoop) GetStats(mailId string) *model.EmailHealthStatus {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if healthStatus, ok := g.stats[mailId]; ok {
		return healthStatus
	}

	return nil
}
