package service

import (
	"fmt"

	"github.com/AbhinitKumarRai/email-health-service/internal/feedback"
	"github.com/AbhinitKumarRai/email-health-service/pkg/model"
)

type Service struct {
	feedbackRegistry  map[feedback.DomainType]feedback.FeedbackLoop
	eventInputChannel chan model.EmailEvent
}

func NewService(eventChan chan model.EmailEvent) *Service {
	service := &Service{
		eventInputChannel: eventChan,
		feedbackRegistry:  make(map[feedback.DomainType]feedback.FeedbackLoop),
	}

	service.registerFeedbackLoops()

	go service.consumeEmailEvents()

	return service
}

func (s *Service) registerFeedbackLoops() {
	s.feedbackRegistry[feedback.Google] = feedback.CreateClient(feedback.Google)
	s.feedbackRegistry[feedback.Yahoo] = feedback.CreateClient(feedback.Yahoo)
}

func (s *Service) consumeEmailEvents() {
	for evt := range s.eventInputChannel {
		for _, loop := range s.feedbackRegistry {
			loop.RegisterMailID(evt.MailID)
		}
	}
}

func (s *Service) GetStats(emailID string) (*model.EmailHealthStatus, error) {
	res := &model.EmailHealthStatus{}

	for _, feedbackInst := range s.feedbackRegistry {
		feedback := feedbackInst.GetStats(emailID)

		if feedback != nil {
			res.EmailID = feedback.EmailID
			res.Delivered += feedback.Delivered
			res.Opened += feedback.Opened
			res.Read += feedback.Read
			res.Spam += feedback.Spam
		}
	}

	fmt.Println(res)
	return res, nil
}

func (s *Service) GetAllMailStats() (map[string]*model.EmailHealthStatus, error) {
	res := make(map[string]*model.EmailHealthStatus)

	for _, feedbackInst := range s.feedbackRegistry {
		feedback := feedbackInst.GetAllStats()

		for mailId, health := range feedback {
			if _, ok := res[mailId]; !ok {
				res[mailId] = health
			} else {
				res[mailId].Delivered += feedback[mailId].Delivered
				res[mailId].Opened += feedback[mailId].Opened
				res[mailId].Read += feedback[mailId].Read
				res[mailId].Spam += feedback[mailId].Spam
			}
		}
	}
	return res, nil
}
