package httpendpoint

import (
	"net/http"
)

func (s *Server) signOutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		if err := s.authServer.SignOut(r.Context(), accessToken); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
