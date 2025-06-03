package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/models"
	"api/internal/server"
	"api/internal/service"
	"api/pkg/cache"
	"api/pkg/database"
	"api/pkg/services/worker"
	"api/pkg/storage"
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

	pgDB, err := database.Init(cfg.LogLevel, cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Username, cfg.Database.Password, cfg.Name)
	if err != nil {
		panic(err)
	}
	err = pgDB.AutoMigrate(models.GetModels()...)
	if err != nil {
		panic(err)
	}

	cache := cache.Init(cfg)

	storageProvider := storage.Init(&storage.Config{
		Provider: cfg.Provider, Root: cfg.DiskRoot,
		Endpoint: cfg.Endpoint, AccessKey: cfg.AccessKey, SecretKey: cfg.SecretKey,
	})

	service := service.Init(cfg, pgDB, storageProvider)
	grpcServer := server.StartGRPCServer(cfg, service)

	err = pgDB.Callback().Query().Register("mediaItemUrl", (&models.MediaItemURLPlugin{
		Storage: storageProvider,
		Cache:   cache,
	}).TransformMediaItemURL)
	if err != nil {
		panic(err)
	}

	handler := &handlers.Handler{
		Config: cfg,
		DB:     pgDB,
		Cache:  cache,
	}
	httpServer := server.StartHTTPServer(handler)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Worker.Host, cfg.Worker.Port), opts...)
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
