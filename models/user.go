package models

import "time"

type User struct {
	Name      string
	Password  string
	CreatedOn time.Time
}
