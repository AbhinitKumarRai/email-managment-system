package feedback

import (
	"math/rand"
	"sync"
	"time"

	"github.com/AbhinitKumarRai/email-health-service/pkg/model"
)

type YahooFeedbackLoop struct {
	stats  map[string]*model.EmailHealthStatus
	ticker *time.Ticker
	mu     sync.RWMutex
}

func NewYahooFeedbackLoop() *YahooFeedbackLoop {
	yahooFeedbackInst := &YahooFeedbackLoop{
		stats:  make(map[string]*model.EmailHealthStatus),
		ticker: time.NewTicker(5 * time.Second),
	}

	go yahooFeedbackInst.simulateFeedback()

	return yahooFeedbackInst
}

func (g *YahooFeedbackLoop) RegisterMailID(emailID string) {
	g.mu.Lock()
	g.stats[emailID] = &model.EmailHealthStatus{
		EmailID: emailID,
	}

	g.mu.Unlock()
}

// Background loop that periodically simulates feedback for all email IDs
func (g *YahooFeedbackLoop) simulateFeedback() {
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

func (g *YahooFeedbackLoop) GetAllStats() map[string]*model.EmailHealthStatus {
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

func (g *YahooFeedbackLoop) GetStats(mailId string) *model.EmailHealthStatus {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if healthStatus, ok := g.stats[mailId]; ok {
		return healthStatus
	}

	return nil
}
