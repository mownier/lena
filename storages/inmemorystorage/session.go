package inmemorystorage

import (
	"context"
	"fmt"
	"lena/errors"
	"lena/models"
)

func (s *InMemoryStorage) AddSession(ctx context.Context, session models.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.sessions[session.AccessToken]
	if exists {
		domain := fmt.Sprintf("inmemorystorage.InMemoryStorage.AddSession: session = %v", session)
		return errors.NewAppError(errors.ErrCodeSessionAlreadyExists, domain, nil)
	}
	s.sessions[session.AccessToken] = session
	return nil
}

func (s *InMemoryStorage) GetSessionByAccessToken(ctx context.Context, accessToken string) (models.Session, error) {
	return s.getSessionByAccessToken(accessToken, true)
}

func (s *InMemoryStorage) UpdateSessionByAccessToken(ctx context.Context, accessToken string, update models.SessionUpdate) (models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, err := s.getSessionByAccessToken(accessToken, false)
	if err != nil {
		domain := fmt.Sprintf("inmemorystorage.InMemoryStorage.UpdateSessionByAccessToken: accessToken = %s, update = %v", accessToken, update)
		return models.Session{}, errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
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

func (s *InMemoryStorage) getSessionByAccessToken(accessToken string, useMutex bool) (models.Session, error) {
	if useMutex {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	session, exists := s.sessions[accessToken]
	if !exists {
		domain := fmt.Sprintf("inmemorystorage.InMemoryStorage.getSessionByAccessToken: accessToken = %v, useMutex = %v", accessToken, useMutex)
		return models.Session{}, errors.NewAppError(errors.ErrCodeSessionDoesNotExist, domain, nil)
	}
	return session, nil
}
