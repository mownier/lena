package models

type RegisterRequest struct {
	Name string `json:"name"`
}

type SignInRequest struct {
	Name string `json:"name"`
}

type SignOutRequest struct {
	AuthenticationKey string `json:"authentication_key"`
}

type VerifyRequest struct {
	AuthenticationKey string `json:"authentication_key"`
}
