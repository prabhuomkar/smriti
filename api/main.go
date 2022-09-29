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

	pgDB, err := database.Init(cfg.Database.LogLevel, cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Name)
	if err != nil {
		panic(err)
	}

	// todo(omkar): initialize GRPC client connection

	service := &service.Service{
		Config: cfg,
		DB:     pgDB,
	}
	server.InitGRPCServer(cfg, service)

	handler := &handlers.Handler{
		Config: cfg,
		DB:     pgDB,
	}
	server.InitHTTPServer(cfg, handler)
	// todo(omkar): handling graceful shutdowns, grpc/http timeouts and reconnections
}
