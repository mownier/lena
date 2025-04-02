package models

import (
	"fmt"
	"time"
)

type User struct {
	Name      string
	Password  string
	CreatedOn time.Time
}

func (u User) String() string {
	return fmt.Sprintf("User{Name: %v, Password: ****, CreatedOn: %v}", u.Name, u.CreatedOn)
}
