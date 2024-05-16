package server

import (
	"api/config"
	"api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartStopGRPCServer(t *testing.T) {
	srv := StartGRPCServer(&config.Config{}, &service.Service{})
	defer srv.Stop()
	assert.NotNil(t, srv)
	StopGRPCServer(srv)
}
