package httpendpoint

import (
	"fmt"
	"lena/auth"
	"lena/config"
	"log"
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

func ListenAndServe(config config.Config, authServer *auth.Server) {
	mux := http.NewServeMux()
	server := NewServer(authServer)
	server.Setup(mux)
	fmt.Printf("lena HTTP server listening on: http://%s:%d\n", config.LocalIP, config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux); err != nil {
		log.Fatalln("failed to serve:", err)
	}
}
