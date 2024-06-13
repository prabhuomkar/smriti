package server

import (
	"api/config"
	"api/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartStopGRPCServer(t *testing.T) {
	cfg := &config.Config{GRPC: config.GRPC{Host: "localhost", Port: 50051}}
	service := &service.Service{Config: cfg}
	srv := StartGRPCServer(cfg, service)
	assert.NotNil(t, srv)
	time.Sleep(time.Second * 5)
	StopGRPCServer(srv)
}
