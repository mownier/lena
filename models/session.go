package models

import (
	"time"
)

type Session struct {
	AccessToken        string
	RefreshToken       string
	UserName           string
	AccesTokenExpiry   time.Time
	RefreshTokenExpiry time.Time
	CreatedOn          time.Time
	ArchivedOn         time.Time
	Archived           bool
}

type SessionUpdate struct {
	ArchivedOn *time.Time
	Archived   *bool
}
