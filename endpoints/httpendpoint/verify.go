package httpendpoint

import (
	"encoding/json"
	"lena/errors"
	"net/http"
)

func (s *Server) verifyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := "httpendpoint.Server.verifyHandler"
		if r.Method != http.MethodPost {
			appError := errors.NewAppError(errors.ErrCodeHTTPMethodNotAllowed, domain, nil)
			response := appError.AsUserFriendlyResponse()
			var message string
			jsonData, jsonErr := json.Marshal(response)
			if jsonErr != nil {
				message = response.Message
			} else {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusMethodNotAllowed)
			return
		}
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			appError := errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, nil)
			response := appError.AsUserFriendlyResponse()
			jsonData, jsonErr := json.Marshal(response)
			message := response.Message
			if jsonErr == nil {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusUnauthorized)
			return
		}
		if err := s.authServer.Verify(r.Context(), accessToken); err != nil {
			appError := errors.NewAppError(errors.ErrCodeVerifyingAccessToken, domain, err)
			response := appError.AsUserFriendlyResponse()
			jsonData, jsonErr := json.Marshal(response)
			message := response.Message
			if jsonErr == nil {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
