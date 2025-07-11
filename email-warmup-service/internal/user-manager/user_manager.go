package usermanager

import (
	"fmt"
	"sync"

	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
)

type UserManager struct {
	Users        map[string]*model.User // emailId to user  map
	GlobalRWLock sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		Users:        make(map[string]*model.User),
		GlobalRWLock: sync.RWMutex{},
	}
}

func (s *UserManager) AddUser(user *model.User) error {
	s.GlobalRWLock.Lock()
	defer s.GlobalRWLock.Unlock()
	if _, ok := s.Users[user.EmailId]; ok {
		return fmt.Errorf("user with emailId: %s is already present", user.EmailId)
	}

	s.Users[user.EmailId] = user

	return nil
}

func (s *UserManager) ListAllUsers() ([]model.User, error) {
	s.GlobalRWLock.RLock()
	defer s.GlobalRWLock.RUnlock()

	var res []model.User
	for _, user := range s.Users {
		res = append(res, *user)
	}
	return res, nil
}

func (s *UserManager) GetUser(emailId string) (model.User, error) {
	if user, ok := s.Users[emailId]; ok {
		return *user, nil
	}

	return model.User{}, fmt.Errorf("user with emailId: %s is already present", emailId)
}
