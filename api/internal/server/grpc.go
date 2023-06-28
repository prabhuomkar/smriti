package server

import (
	"api/config"
	"api/internal/service"
	"api/pkg/services/api"
	"fmt"
	"log"
	"net"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// StartGRPCServer ...
func StartGRPCServer(cfg *config.Config, service *service.Service) *grpc.Server {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		panic(err)
	}

	grpcMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	prometheus.DefaultRegisterer.MustRegister(grpcMetrics)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcMetrics.UnaryServerInterceptor(),
		),
	)
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
