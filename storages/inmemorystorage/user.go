package inmemorystorage

import (
	"context"
	"errors"
	"lena/models"
)

func (s *InMemoryStorage) AddUser(ctx context.Context, user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.Name]; exists {
		return errors.New("user already exists")
	}
	s.users[user.Name] = user
	return nil
}

func (s *InMemoryStorage) GetUserByName(ctx context.Context, name string) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.users[name]
	if !exists {
		return models.User{}, errors.New("user does not exist")
	}
	return user, nil
}
