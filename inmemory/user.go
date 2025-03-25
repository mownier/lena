package inmemory

import (
	"errors"
	"lena/models"
)

var users = make(map[string]models.User)

func AddUser(user models.User) error {
	_, exists := users[user.Name]
	if exists {
		return errors.New("user already exists")
	}
	users[user.Name] = user
	return nil
}

func UserDoesExist(name string) bool {
	_, exists := users[name]
	return exists
}

func UserDoesNotExist(name string) bool {
	return !UserDoesExist(name)
}

func UserDoesHaveAuthentication(name string) bool {
	list, exists := userNameAuthenticationKeys[name]
	if !exists || len(list) == 0 {
		return false
	}
	return true
}
