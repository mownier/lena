package main

import (
	"fmt"
	endpoints "lena/jsonendpoints"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/signin", endpoints.SignIn)
	http.HandleFunc("/register", endpoints.Register)
	http.HandleFunc("/signout", endpoints.SignOut)
	http.HandleFunc("/verify", endpoints.Verify)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Error getting network interfaces: %v", err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP := ipnet.IP.String()
				fmt.Printf("lena listening on: http://%s:%s\n", localIP, port)
			}
		}
	}
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("Error startng server:", err)
	}
}
