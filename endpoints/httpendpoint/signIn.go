package httpendpoint

import (
	"encoding/json"
	"io"
	"lena/errors"
	"net/http"
)

func (s *Server) signInHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := "httpendpoint.Server.signInHandler"
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			appError := errors.NewAppError(errors.ErrCodeHTTPBodyCannotBeRead, domain, err)
			response := appError.AsUserFriendlyResponse()
			var message string
			jsonData, jsonErr := json.Marshal(response)
			if jsonErr != nil {
				message = response.Message
			} else {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var request SignInRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			appError := errors.NewAppError(errors.ErrCodeHTTPBodyMalformed, domain, err)
			response := appError.AsUserFriendlyResponse()
			var message string
			jsonData, jsonErr := json.Marshal(response)
			if jsonErr != nil {
				message = response.Message
			} else {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		session, err := s.authServer.SignIn(r.Context(), request.Name, request.Password)
		if err != nil {
			appError := errors.NewAppError(errors.ErrCodeSigningIn, domain, err)
			var response errors.UserFriendlyResponse
			if other, contains := appError.ContainsCode(errors.ErrCodeUserDoesNotExist); contains {
				response = other.AsUserFriendlyResponse()
			} else if other, contains := appError.ContainsCode(errors.ErrCodeInvalidPassword); contains {
				response = other.AsUserFriendlyResponse()
			} else {
				response = appError.AsUserFriendlyResponse()
			}
			var message string
			jsonData, jsonErr := json.Marshal(response)
			if jsonErr != nil {
				message = response.Message
			} else {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		response := SignInResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			ExpiresOn:    session.AccesTokenExpiry,
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			appError := errors.NewAppError(errors.ErrCodeGeneratingResponse, domain, err)
			response := appError.AsUserFriendlyResponse()
			var message string
			jsonData, jsonErr := json.Marshal(response)
			if jsonErr != nil {
				message = response.Message
			} else {
				message = string(jsonData)
			}
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
