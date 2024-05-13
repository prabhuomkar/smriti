package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/jobs"
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

//nolint:funlen
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

	cache := cache.Init(cfg)

	storageProvider := storage.Init(&storage.Config{
		Provider: cfg.Storage.Provider, Root: cfg.Storage.DiskRoot,
		Endpoint: cfg.Storage.Endpoint, AccessKey: cfg.Storage.AccessKey, SecretKey: cfg.Storage.SecretKey,
	})

	service := &service.Service{
		Config:  cfg,
		DB:      pgDB,
		Storage: storageProvider,
	}
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
		grpc.WithBlock(),
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Worker.Host, cfg.Worker.Port), opts...)
	if err != nil {
		panic(err)
	}
	handler.Worker = worker.NewWorkerClient(conn)

	jobsInstance := &jobs.Job{
		Config:  cfg,
		DB:      pgDB,
		Storage: storageProvider,
		Worker:  worker.NewWorkerClient(conn),
	}
	go jobsInstance.StartJobs()

	// graceful shutdown
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal
	server.StopGRPCServer(grpcServer)
	server.StopHTTPServer(httpServer)
}
