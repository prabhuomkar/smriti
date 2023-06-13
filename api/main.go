package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/models"
	"api/internal/server"
	"api/internal/service"
	"api/internal/storage"
	"api/pkg/cache"
	"api/pkg/database"
	"api/pkg/services/worker"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "go.uber.org/automaxprocs"
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
	err = pgDB.AutoMigrate(models.GetModels()...)
	if err != nil {
		panic(err)
	}

	service := &service.Service{
		Config:  cfg,
		DB:      pgDB,
		Storage: storage.Init(cfg.Storage),
	}
	grpcServer := server.StartGRPCServer(cfg, service)

	cache := cache.Init(cfg)

	handler := &handlers.Handler{
		Config: cfg,
		DB:     pgDB,
		Cache:  cache,
	}
	httpServer := server.StartHTTPServer(handler)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Worker.Host, cfg.Worker.Port), opts...)
	if err != nil {
		panic(err)
	}
	handler.Worker = worker.NewWorkerClient(conn)

	// graceful shutdown
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal
	server.StopGRPCServer(grpcServer)
	server.StopHTTPServer(httpServer)
}
