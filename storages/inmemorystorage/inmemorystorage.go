package inmemorystorage

import (
	"lena/models"
	"sync"
)

type InMemoryStorage struct {
	mu       sync.RWMutex
	users    map[string]models.User
	sessions map[string]models.Session
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:    make(map[string]models.User),
		sessions: make(map[string]models.Session),
	}
}
