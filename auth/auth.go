package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"lena/models"
	"lena/storages"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	storage storages.Storage
}

func NewServer(storage storages.Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) Register(ctx context.Context, name string, password string) (models.Session, error) {
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return models.Session{}, err
	}
	user := s.newUser(name, hashedPassword)
	err = s.storage.AddUser(ctx, user)
	if err != nil {
		return models.Session{}, err
	}
	session := s.newSession(name)
	err = s.storage.AddSession(ctx, session)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (s *Server) SignIn(ctx context.Context, name string, password string) (models.Session, error) {
	user, err := s.storage.GetUserByName(ctx, name)
	if err != nil {
		return models.Session{}, err
	}
	err = s.verifyPassword(password, user.Password)
	if err != nil {
		return models.Session{}, err
	}
	session := s.newSession(name)
	err = s.storage.AddSession(ctx, session)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (s *Server) SignOut(ctx context.Context, accessToken string) error {
	session, err := s.storage.GetSessionByAccessToken(ctx, accessToken)
	if err != nil {
		return nil
	}
	if session.Archived {
		return errors.New("session is already invalidated")
	}
	_, err = s.storage.UpdateSessionByAccessToken(ctx, accessToken, s.newSessionUpdateForArchiving())
	if err != nil {
		return nil
	}
	return nil
}

func (s *Server) Verify(ctx context.Context, accessToken string) error {
	if _, err := s.verify(ctx, accessToken); err != nil {
		return err
	}
	return nil
}

func (s *Server) Refresh(ctx context.Context, accessToken string, refreshToken string) (models.Session, error) {
	session, err := s.verify(ctx, accessToken)
	if err != nil {
		return models.Session{}, err
	}
	if session.RefreshToken != refreshToken {
		return models.Session{}, errors.New("invalid refresh token")
	}
	_, err = s.storage.UpdateSessionByAccessToken(ctx, accessToken, s.newSessionUpdateForArchiving())
	if err != nil {
		return models.Session{}, err
	}
	session = s.newSession(session.UserName)
	err = s.storage.AddSession(ctx, session)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (s *Server) verify(ctx context.Context, accessToken string) (models.Session, error) {
	session, err := s.storage.GetSessionByAccessToken(ctx, accessToken)
	if err != nil {
		return models.Session{}, err
	}
	if session.Archived {
		return models.Session{}, errors.New("session is already invalidated")
	}
	now := time.Now().UTC()
	accesTokenExpired := now.Equal(session.AccesTokenExpiry) || now.After(session.AccesTokenExpiry)
	refreshTokenExpired := now.Equal(session.RefreshTokenExpiry) || now.After(session.RefreshTokenExpiry)
	if accesTokenExpired {
		return models.Session{}, errors.New("session is already expired")
	}
	if refreshTokenExpired {
		s.storage.UpdateSessionByAccessToken(ctx, accessToken, s.newSessionUpdateForArchiving())
		return models.Session{}, errors.New("session can no longer be extended")
	}
	return session, nil
}

func (s *Server) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Server) verifyPassword(input string, stored string) error {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(input))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func (s *Server) newSessionUpdateForArchiving() models.SessionUpdate {
	archivedOn := time.Now().UTC()
	archived := true
	return models.SessionUpdate{
		ArchivedOn: &archivedOn,
		Archived:   &archived,
	}
}

func (s *Server) newUser(name string, password string) models.User {
	return models.User{
		Name:      name,
		Password:  password,
		CreatedOn: time.Now().UTC(),
	}
}

func (s *Server) newToken(variable string) string {
	random := rand.Intn(1_000_000_000_000)
	now := time.Now().UnixNano()
	input := fmt.Sprintf("%s-%d-%d", variable, now, random)
	hash := sha256.New()
	hash.Write([]byte(input))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func (s *Server) newSession(userName string) models.Session {
	accessToken := s.newToken(userName)
	return models.Session{
		AccessToken:        accessToken,
		RefreshToken:       s.newToken(accessToken),
		UserName:           userName,
		AccesTokenExpiry:   time.Now().AddDate(0, 0, 7).UTC(),  // 7 days
		RefreshTokenExpiry: time.Now().AddDate(0, 0, 30).UTC(), // 30 days
		CreatedOn:          time.Now().UTC(),
		ArchivedOn:         time.Now().AddDate(1000, 0, 0).UTC(),
		Archived:           false,
	}
}
