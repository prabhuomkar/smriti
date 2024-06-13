package server

import (
	"api/config"
	"api/internal/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartStopHTTPServer(t *testing.T) {
	handler := &handlers.Handler{Config: &config.Config{Storage: config.Storage{Provider: "disk"}}}
	srv := StartHTTPServer(handler)
	defer srv.Close()
	assert.NotNil(t, srv)
	StopHTTPServer(srv)
}

func TestGetMiddlewareFuncs(t *testing.T) {
	mockConfig := &config.Config{}
	tests := []struct {
		Name        string
		JWTCheck    bool
		Features    []string
		ExpectedLen int
	}{
		{
			Name:        "without jwt check and no features",
			JWTCheck:    false,
			Features:    []string{},
			ExpectedLen: 0,
		},
		{
			Name:        "without jwt check and features",
			JWTCheck:    false,
			Features:    []string{"places", "favourites"},
			ExpectedLen: 2,
		},
		{
			Name:        "with jwt check and no features",
			JWTCheck:    true,
			Features:    []string{},
			ExpectedLen: 1,
		},
		{
			Name:        "with jwt check and features",
			JWTCheck:    true,
			Features:    []string{"places", "favourites"},
			ExpectedLen: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			middlewareFuncs := getMiddlewareFuncs(mockConfig, nil, tc.JWTCheck, tc.Features...)
			assert.Equal(t, tc.ExpectedLen, len(middlewareFuncs))
		})
	}
}
