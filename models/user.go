package models

import "time"

type User struct {
	Name      string
	CreatedOn time.Time
}

func GenerateUser(name string) User {
	return User{
		Name:      name,
		CreatedOn: time.Now().UTC(),
	}
}
