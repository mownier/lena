package main

import (
	"errors"
	"fmt"
	"lena/auth"
	"lena/endpoints/grpcendpoint"
	"lena/endpoints/httpendpoint"
	"lena/storages"
	"lena/storages/inmemorystorage"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln("Error getting network interfaces:", err)
	}
	localIP := ""
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	port, err := getPort()
	if err != nil {
		log.Fatalln("Error getting port:", err)
	}
	storage, err := getStorage()
	if err != nil {
		log.Fatalln("Error getting storage:", err)
	}
	endpoint, err := getEndpoint()
	if err != nil {
		log.Fatalln("Error getting storage:", err)
	}
	enableReflection, err := getReflection()
	if err != nil {
		log.Fatalln("Error getting reflection:", err)
	}
	var userStorage storages.UserStorage
	var sessionStorage storages.SessionStorage
	if storage == "inmemory" {
		inmemorystorage := inmemorystorage.NewInMemoryStorage()
		userStorage = inmemorystorage
		sessionStorage = inmemorystorage
	}
	authServer := auth.NewServer(userStorage, sessionStorage)
	if endpoint == "http" {
		runHTTPEndpoint(localIP, port, authServer)
		return
	}
	if endpoint == "grpc" {
		runGRPCEndpoint(localIP, port, enableReflection, authServer)
		return
	}
}

func getPort() (int, error) {
	input := os.Getenv("LENA_PORT")
	if input == "" {
		return 8080, nil
	}
	port, err := strconv.Atoi(input)
	if err != nil {
		return -1, err
	}
	return port, nil
}

func getStorage() (string, error) {
	array := []string{"inmemory", "sqlite"}
	input := os.Getenv("LENA_STORAGE")
	if input == "" {
		return array[0], nil
	}
	input = strings.ToLower(input)
	for _, value := range array {
		if input == value {
			return input, nil
		}
	}
	return "", errors.New("selected storage option not found")
}

func getEndpoint() (string, error) {
	array := []string{"http", "grpc"}
	input := os.Getenv("LENA_ENDPOINT")
	if input == "" {
		return array[0], nil
	}
	input = strings.ToLower(input)
	for _, value := range array {
		if input == value {
			return input, nil
		}
	}
	return "", errors.New("selected endpoint option not found")
}

func getReflection() (bool, error) {
	input := os.Getenv("LENA_REFLECTION")
	if input == "" {
		return false, nil
	}
	input = strings.ToLower(input)
	if input == "false" {
		return false, nil
	}
	if input == "true" {
		return true, nil
	}
	return false, errors.New("selected reflection option not found")
}

func runHTTPEndpoint(localIP string, port int, authServer *auth.Server) {
	mux := http.NewServeMux()
	server := httpendpoint.NewServer(authServer)
	server.Setup(mux)
	fmt.Printf("lena HTTP server listening on: http://%s:%d\n", localIP, port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalln("failed to serve:", err)
	}
}

func runGRPCEndpoint(localIP string, port int, enableReflection bool, authServer *auth.Server) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	if enableReflection {
		reflection.Register(grpcServer)
	}
	grpcendpoint.RegisterLenaServiceServer(grpcServer, grpcendpoint.NewServer(authServer))
	fmt.Printf("lena GRPC server listening on: tcp://%s:%d\n", localIP, port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln("failed to server:", err)
	}
}
