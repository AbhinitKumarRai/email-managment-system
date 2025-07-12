package feedback

import (
	"math/rand"
	"time"
)

type GmailFeedbackLoop struct{}

func (g *GmailFeedbackLoop) Check(emailID int64, recipient string) (delivered bool, spam bool, feedbacks int, err error) {
	rand.Seed(time.Now().UnixNano() + emailID)
	feedbacks = 1
	if rand.Intn(10) < 7 {
		return true, false, feedbacks, nil // 70% delivered
	}
	return false, true, feedbacks, nil // 30% spam
}

func init() {
	Register("gmail.com", &GmailFeedbackLoop{})
}
