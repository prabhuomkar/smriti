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

// InitGRPCServer ...
func InitGRPCServer(cfg *config.Config, service *service.Service) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	api.RegisterAPIServiceServer(server, service)

	go func() {
		log.Printf("starting grpc server on: %d", cfg.GRPC.Port)
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()
}
