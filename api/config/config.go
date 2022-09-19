package config

import "github.com/kelseyhightower/envconfig"

type (
	// Log ...
	Log struct {
		Level string `default:"INFO"`
	}

	// API ...
	API struct {
		Host string `default:"127.0.0.1"`
		Port int    `default:"5001"`
	}

	// GRPC ...
	GRPC struct {
		Host string `default:"127.0.0.1"`
		Port int    `default:"15001"`
	}

	// Database ...
	Database struct {
		Host     string `default:"db"`
		Port     int    `default:"5432"`
		Username string `default:"pensieve"`
		Password string `default:"pensieve"`
	}

	// Feature ...
	Feature struct {
		Favourites    bool `default:"true"`
		Hidden        bool `default:"true"`
		Trash         bool `default:"true"`
		Albums        bool `default:"true"`
		Explore       bool `default:"true"`
		ExplorePlaces bool `default:"true"`
		ExploreThings bool `default:"true"`
		ExplorePeople bool `default:"true"`
		Sharing       bool `default:"true"`
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
