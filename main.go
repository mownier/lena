package main

import (
	"lena/auth"
	"lena/config"
	"lena/endpoints/grpcendpoint"
	"lena/endpoints/httpendpoint"
	"lena/storages"
	"lena/storages/inmemorystorage"
	"log"
)

func main() {
	config, err := config.Setup()
	if err != nil {
		log.Fatalln("Error getting config:", err)
	}
	var store storages.Storage
	if config.Storage == "inmemory" {
		store = inmemorystorage.NewInMemoryStorage()
	}
	authServer := auth.NewServer(store)
	if config.Endpoint == "http" {
		httpendpoint.ListenAndServe(config, authServer)
		return
	}
	if config.Endpoint == "grpc" {
		grpcendpoint.ListenAndServe(config, authServer)
		return
	}
}
