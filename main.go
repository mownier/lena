package main

import (
	"fmt"
	"lena/endpoints/httpendpoint"
	"lena/server"
	"lena/storages/inmemorystorage"
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
	storage := inmemorystorage.NewInMemoryStorage()
	authServer := server.NewAuthServer(storage, storage)
	mux := http.NewServeMux()
	httpendpoint.SetupHTTPHandlers(mux, authServer)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalln("Error startng server:", err)
	}
}
