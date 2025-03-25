package inmemory

import (
	"errors"
	"lena/models"
	"time"
)

var authentications = make(map[string]models.Authentication)
var userNameAuthenticationKeys = make(map[string][]string)

func AddAuthentication(authentication models.Authentication) error {
	_, exists := authentications[authentication.Key]
	if exists {
		return errors.New("authentication does already exist")
	}
	authentications[authentication.Key] = authentication
	_, exists = userNameAuthenticationKeys[authentication.UserName]
	if !exists {
		userNameAuthenticationKeys[authentication.UserName] = []string{}
	}
	userNameAuthenticationKeys[authentication.UserName] = append(
		userNameAuthenticationKeys[authentication.UserName],
		authentication.Key,
	)
	return nil
}

func AuthenticationDoesExist(key string) bool {
	_, exists := authentications[key]
	return exists
}

func AuthenticationDoesNotExist(key string) bool {
	return !AuthenticationDoesExist(key)
}

func ArchiveAuthentication(key string) error {
	authentication, exists := authentications[key]
	if !exists {
		return errors.New("authentication does not exist")
	}
	authentication.Archived = true
	authentication.ArchivedOn = time.Now().UTC()
	authentications[key] = authentication
	list, exists := userNameAuthenticationKeys[authentication.UserName]
	if !exists {
		return nil
	}
	updatedList := []string{}
	for _, authenticationKey := range list {
		if authenticationKey != authentication.Key {
			updatedList = append(updatedList, authenticationKey)
		}
	}
	userNameAuthenticationKeys[authentication.UserName] = updatedList
	return nil
}

func AuthenticationIsArchived(key string) (bool, error) {
	authentication, exists := authentications[key]
	if !exists {
		return false, errors.New("authentication does not exist")
	}
	return authentication.Archived, nil
}
