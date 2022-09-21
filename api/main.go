package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/server"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	// todo(omkar): initialize DB connection
	// todo(omkar): initialize GRPC client connection

	log.Printf("starting grpc server on: %d", cfg.API.Port)
	err = server.InitGRPCServer(cfg)
	if err != nil {
		panic(err)
	}

	handler := &handlers.Handler{
		Config: cfg,
	}

	log.Printf("starting http api server on: %d", cfg.API.Port)
	err = server.InitHTTPServer(cfg, handler)
	if err != http.ErrServerClosed {
		panic(err)
	}
}
