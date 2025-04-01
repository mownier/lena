package inmemorystorage

import (
	"context"
	"errors"
	"lena/models"
)

func (s *InMemoryStorage) AddSession(ctx context.Context, session models.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.sessions[session.AccessToken]
	if exists {
		return errors.New("session already exists")
	}
	s.sessions[session.AccessToken] = session
	return nil
}

func (s *InMemoryStorage) GetSessionByAccessToken(ctx context.Context, accessToken string) (models.Session, error) {
	return s.getSessionByAccessToken(ctx, accessToken, true)
}

func (s *InMemoryStorage) UpdateSessionByAccessToken(ctx context.Context, accessToken string, update models.SessionUpdate) (models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, err := s.getSessionByAccessToken(ctx, accessToken, false)
	if err != nil {
		return models.Session{}, err
	}
	if update.ArchivedOn != nil {
		session.ArchivedOn = *update.ArchivedOn
	}
	if update.Archived != nil {
		session.Archived = *update.Archived
	}
	s.sessions[accessToken] = session
	return session, nil
}

func (s *InMemoryStorage) getSessionByAccessToken(ctx context.Context, accessToken string, useMutex bool) (models.Session, error) {
	if useMutex {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	session, exists := s.sessions[accessToken]
	if !exists {
		return models.Session{}, errors.New("session does not exist")
	}
	return session, nil
}
