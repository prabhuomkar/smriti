package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/server"
	"api/internal/service"
	"api/pkg/cache"
	"api/pkg/database"
	"api/pkg/services/worker"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Worker.Host, cfg.Worker.Port), opts...)
	if err != nil {
		panic(err)
	}
	workerClient := worker.NewWorkerClient(conn)

	service := &service.Service{
		Config: cfg,
		DB:     pgDB,
	}
	server.InitGRPCServer(cfg, service)

	cache, err := cache.Init()
	if err != nil {
		panic(err)
	}

	handler := &handlers.Handler{
		Config: cfg,
		DB:     pgDB,
		Worker: &workerClient,
		Cache:  cache,
	}
	server.InitHTTPServer(cfg, handler)
	// todo(omkar): handling graceful shutdowns, grpc/http timeouts and reconnections
}
