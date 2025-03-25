package models

type RegisterResponse struct {
	Okay              bool   `json:"okay"`
	ErrorMessage      string `json:"error_message"`
	ErrorCode         int    `json:"error_code"`
	AuthenticationKey string `json:"authentication_key"`
}

type SignInResponse struct {
	Okay              bool   `json:"okay"`
	ErrorMessage      string `json:"error_message"`
	ErrorCode         int    `json:"error_code"`
	AuthenticationKey string `json:"authentication_key"`
}

type SignOutResponse struct {
	Okay         bool   `json:"okay"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    int    `json:"error_code"`
}

type VerifyResponse struct {
	Okay         bool   `json:"okay"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    int    `json:"error_code"`
}
