package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/server"
	"api/internal/service"
	"api/pkg/database"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	db, err := database.Init(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username,
		cfg.Database.Password, cfg.Database.Name)
	if err != nil {
		panic(err)
	}

	// todo(omkar): initialize GRPC client connection

	service := &service.Service{
		Config: cfg,
		DB:     db,
	}
	server.InitGRPCServer(cfg, service)

	handler := &handlers.Handler{
		Config: cfg,
		DB:     db,
	}
	server.InitHTTPServer(cfg, handler)
	// todo(omkar): handling graceful shutdowns, grpc/http timeouts and reconnections
}
