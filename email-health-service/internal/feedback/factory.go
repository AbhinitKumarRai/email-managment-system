package feedback

func CreateClient(domain DomainType) FeedbackLoop {
	switch domain {
	case Google:
		return NewGmailFeedbackLoop()
	case Yahoo:
		return NewYahooFeedbackLoop()
	}

	return nil
}
