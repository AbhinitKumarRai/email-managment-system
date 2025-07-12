package usermanager

import (
	"fmt"
	"sync"
	"time"

	"github.com/AbhinitKumarRai/email-warmup-service/pkg/model"
)

type UserManager struct {
	users        map[string]*model.User
	emailIdList  []string
	globalRWLock sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		users:        make(map[string]*model.User),
		emailIdList:  []string{},
		globalRWLock: sync.RWMutex{},
	}
}

func (s *UserManager) AddUser(user *model.User) error {
	s.globalRWLock.Lock()
	defer s.globalRWLock.Unlock()
	if _, ok := s.users[user.EmailId]; ok {
		return fmt.Errorf("user with emailId: %s is already present", user.EmailId)
	}

	user.CreatedAt = time.Now()
	s.emailIdList = append(s.emailIdList, user.EmailId)
	s.users[user.EmailId] = user

	return nil
}

func (s *UserManager) GetAllUsers() ([]model.User, error) {
	s.globalRWLock.RLock()
	defer s.globalRWLock.RUnlock()

	var res []model.User
	for _, user := range s.users {
		res = append(res, *user)
	}
	return res, nil
}

func (s *UserManager) GetAllEmailIds() ([]string, error) {
	return s.emailIdList, nil
}

func (s *UserManager) GetUser(emailId string) (model.User, error) {
	if user, ok := s.users[emailId]; ok {
		return *user, nil
	}

	return model.User{}, fmt.Errorf("user with emailId: %s is already present", emailId)
}

func (s *UserManager) DeleteUser(emailId string) error {
	s.globalRWLock.Lock()
	defer s.globalRWLock.Unlock()

	if _, ok := s.users[emailId]; ok {
		delete(s.users, emailId)

		for i, email := range s.emailIdList {
			if email == emailId {
				s.emailIdList = append(s.emailIdList[:i], s.emailIdList[i+1:]...)
				return nil
			}
		}
	}

	return nil
}
