package server

import (
	"api/config"
	"api/internal/service"
	"api/pkg/services/api"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// StartGRPCServer ...
func StartGRPCServer(cfg *config.Config, service *service.Service) *grpc.Server {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterAPIServer(grpcServer, service)

	go func() {
		log.Printf("starting grpc server on: %d", cfg.GRPC.Port)
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return grpcServer
}

// StopGRPCServer ...
func StopGRPCServer(grpcServer *grpc.Server) {
	log.Println("stopping grpc server")
	grpcServer.GracefulStop()
}
