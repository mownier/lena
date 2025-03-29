package storages

import (
	"context"
	"lena/models"
)

type Storage interface {
	UserStorage
	SessionStorage
}

type UserStorage interface {
	AddUser(ctx context.Context, user models.User) error
	GetUserByName(ctx context.Context, name string) (models.User, error)
}

type SessionStorage interface {
	AddSession(ctx context.Context, session models.Session) error
	GetSessionByAccessToken(ctx context.Context, accessToken string) (models.Session, error)
	UpdateSessionByAccessToken(ctx context.Context, accessToken string, update models.SessionUpdate) (models.Session, error)
}
