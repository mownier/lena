package httpendpoint

import (
	"lena/server"
	"net/http"
)

func SetupHTTPHandlers(mux *http.ServeMux, authServer *server.AuthServer) {
	mux.HandleFunc("/register", registerHandler(authServer))
	mux.HandleFunc("/signin", signInHandler(authServer))
	mux.HandleFunc("/signout", signOutHandler(authServer))
	mux.HandleFunc("/verify", verifyHandler(authServer))
	mux.HandleFunc("/refresh", refreshHandler(authServer))
}
