package httpendpoint

import (
	"lena/server"
	"net/http"
)

func verifyHandler(server *server.AuthServer) http.HandlerFunc {
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
		if err := server.Verify(r.Context(), accessToken); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
