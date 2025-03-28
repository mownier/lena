package httpendpoint

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
