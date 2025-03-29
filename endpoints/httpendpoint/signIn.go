package httpendpoint

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) signInHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var request SignInRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		session, err := s.authServer.SignIn(r.Context(), request.Name, request.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response := SignInResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			ExpiresOn:    session.AccesTokenExpiry,
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
