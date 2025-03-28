package httpendpoint

import "time"

type RegisterResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresOn    time.Time `json:"expires_on"`
}

type SignInResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresOn    time.Time `json:"expires_on"`
}

type RefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresOn    time.Time `json:"expires_on"`
}
