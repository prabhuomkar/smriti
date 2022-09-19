package config

import "github.com/kelseyhightower/envconfig"

type (
	// Log ...
	Log struct {
		Level string `default:"INFO"`
	}

	// API ...
	API struct {
	}

	// GRPC ...
	GRPC struct {
	}

	// Database ...
	Database struct {
	}

	// Feature ...
	Feature struct {
	}

	// Config ...
	Config struct {
		Log
		API
		GRPC
		Database
		Feature
	}
)

// Init ...
func Init() (*Config, error) {
	var cfg Config
	err := envconfig.Process("pensieve", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
