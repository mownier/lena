package httpendpoint

import (
	"lena/auth"
	"net/http"
)

func SetupHTTPHandlers(mux *http.ServeMux, authServer *auth.Server) {
	mux.HandleFunc("/register", registerHandler(authServer))
	mux.HandleFunc("/signin", signInHandler(authServer))
	mux.HandleFunc("/signout", signOutHandler(authServer))
	mux.HandleFunc("/verify", verifyHandler(authServer))
	mux.HandleFunc("/refresh", refreshHandler(authServer))
}
