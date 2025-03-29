package httpendpoint

import (
	"lena/auth"
	"net/http"
)

type Server struct {
	authServer *auth.Server
}

func NewServer(authServer *auth.Server) *Server {
	return &Server{authServer: authServer}
}

func (s *Server) Setup(mux *http.ServeMux) {
	mux.HandleFunc("/register", s.registerHandler())
	mux.HandleFunc("/signin", s.signInHandler())
	mux.HandleFunc("/signout", s.signOutHandler())
	mux.HandleFunc("/verify", s.verifyHandler())
	mux.HandleFunc("/refresh", s.refreshHandler())
}
