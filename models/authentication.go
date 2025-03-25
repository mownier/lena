package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Authentication struct {
	Key        string
	UserName   string
	CreatedOn  time.Time
	ArchivedOn time.Time
	Archived   bool
}

func GenerateAuthenticationKey(userName string, random int) string {
	now := time.Now().UnixNano()
	input := fmt.Sprintf("%s-%d-%d", userName, now, random)
	hash := sha256.New()
	hash.Write([]byte(input))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func GenerateAuthentication(key string, userName string) Authentication {
	return Authentication{
		Key:        key,
		UserName:   userName,
		CreatedOn:  time.Now().UTC(),
		ArchivedOn: time.Now().AddDate(1000, 0, 0).UTC(),
		Archived:   false,
	}
}
