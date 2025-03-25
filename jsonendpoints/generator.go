package jsonendpoints

import (
	storage "lena/inmemory"
	"lena/models"
	"math/rand"
)

func generateAuthentication(userName string) models.Authentication {
	generatedKey := ""
	for len(generatedKey) == 0 {
		random := rand.Intn(1_000_000_000_000)
		tempKey := models.GenerateAuthenticationKey(userName, random)
		if storage.AuthenticationDoesNotExist(tempKey) {
			generatedKey = tempKey
		}
	}
	return models.GenerateAuthentication(generatedKey, userName)
}
