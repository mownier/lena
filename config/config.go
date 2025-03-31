package config

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	LocalIP    string
	Port       int
	Reflection bool
	Storage    string
	Endpoint   string
}

func Setup() (Config, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return Config{}, err
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
		return Config{}, err
	}
	storage, err := getStorage()
	if err != nil {
		return Config{}, err
	}
	endpoint, err := getEndpoint()
	if err != nil {
		return Config{}, err
	}
	reflection, err := getReflection()
	if err != nil {
		return Config{}, err
	}
	config := Config{
		LocalIP:    localIP,
		Port:       port,
		Storage:    storage,
		Endpoint:   endpoint,
		Reflection: reflection,
	}
	return config, err
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
