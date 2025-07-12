package feedback

type FeedbackLoop interface {
	Check(emailID int64, recipient string) (delivered bool, spam bool, feedbacks int, err error)
}

var registry = make(map[string]FeedbackLoop)

func Register(domain string, loop FeedbackLoop) {
	registry[domain] = loop
}

func Get(domain string) (FeedbackLoop, bool) {
	loop, ok := registry[domain]
	return loop, ok
}
