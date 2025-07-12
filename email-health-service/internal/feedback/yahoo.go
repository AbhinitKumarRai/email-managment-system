package feedback

import (
	"math/rand"
	"time"
)

type YahooFeedbackLoop struct{}

func (y *YahooFeedbackLoop) Check(emailID int64, recipient string) (delivered bool, spam bool, feedbacks int, err error) {
	rand.Seed(time.Now().UnixNano() + emailID)
	feedbacks = 1
	if rand.Intn(10) < 6 {
		return true, false, feedbacks, nil // 60% delivered
	}
	return false, true, feedbacks, nil // 40% spam
}

func init() {
	Register("yahoo.com", &YahooFeedbackLoop{})
}
