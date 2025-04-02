package inmemorystorage

import (
	"context"
	"fmt"
	"lena/errors"
	"lena/models"
)

func (s *InMemoryStorage) AddUser(ctx context.Context, user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.Name]; exists {
		domain := fmt.Sprintf("inmemorystorage.InMemoryStorage.AddUser: user = %v", user)
		return errors.NewAppError(errors.ErrCodeUserAlreadyExists, domain, nil)
	}
	s.users[user.Name] = user
	return nil
}

func (s *InMemoryStorage) GetUserByName(ctx context.Context, name string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.users[name]
	if !exists {
		domain := fmt.Sprintf("inmemorystorage.InMemoryStorage.GetUserByName: name = %v", name)
		return models.User{}, errors.NewAppError(errors.ErrCodeUserDoesNotExist, domain, nil)
	}
	return user, nil
}
