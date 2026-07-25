package main

import (
	"api/config"
	"api/internal/handlers"
	"api/internal/server"
	"api/internal/service"
	"api/pkg/cache"
	"api/pkg/database"
	"api/pkg/storage"
	"os"
	"os/signal"
	"syscall"

	_ "go.uber.org/automaxprocs"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	pgDB, err := database.Init(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username,
		cfg.Database.Password, cfg.Name, cfg.Timeout)
	if err != nil {
		panic(err)
	}

	cache := cache.Init(cfg.Type, cfg.Cache.Host, cfg.Cache.Port, cfg.Cache.Password)

	storageProvider := storage.Init(&storage.Config{
		Provider: cfg.Provider, Root: cfg.DiskRoot, Endpoint: cfg.Endpoint,
		AccessKey: cfg.AccessKey, SecretKey: cfg.SecretKey,
	})

	service := service.Init(cfg, pgDB, storageProvider)
	grpcServer := server.StartGRPCServer(cfg, service)

	handler := &handlers.Handler{
		Config: cfg, DB: pgDB, Cache: cache,
	}
	httpServer := server.StartHTTPServer(handler)

	// graceful shutdown
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal
	server.StopGRPCServer(grpcServer)
	server.StopHTTPServer(httpServer)
}
