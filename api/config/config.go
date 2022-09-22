package config

import "github.com/kelseyhightower/envconfig"

type (
	// Log ...
	Log struct {
		Level string `envconfig:"PENSIEVE_LOG_LEVEL" default:"INFO"`
	}

	// API ...
	API struct {
		Host string `envconfig:"PENSIEVE_API_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"PENSIEVE_API_PORT" default:"5001"`
	}

	// GRPC ...
	GRPC struct {
		Host string `envconfig:"PENSIEVE_GRPC_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"PENSIEVE_GRPC_PORT" default:"15001"`
	}

	// Database ...
	Database struct {
		Host     string `envconfig:"PENSIEVE_DATABASE_HOST" default:"db"`
		Port     int    `envconfig:"PENSIEVE_DATABASE_PORT" default:"5432"`
		Username string `envconfig:"PENSIEVE_DATABASE_USERNAME" default:"pensieve"`
		Password string `envconfig:"PENSIEVE_DATABASE_PASSWORD" default:"pensieve"`
		Name     string `envconfig:"PENSIEVE_DATABASE_NAME" default:"pensieve"`
	}

	// Feature ...
	Feature struct {
		Favourites    bool `envconfig:"PENSIEVE_FEATURE_FAVOURITES" default:"false"`
		Hidden        bool `envconfig:"PENSIEVE_FEATURE_HIDDEN" default:"false"`
		Trash         bool `envconfig:"PENSIEVE_FEATURE_TRASH" default:"false"`
		Albums        bool `envconfig:"PENSIEVE_FEATURE_ALBUMS" default:"false"`
		Explore       bool `envconfig:"PENSIEVE_FEATURE_EXPLORE" default:"false"`
		ExplorePlaces bool `envconfig:"PENSIEVE_FEATURE_EXPLORE_PLACES" default:"false"`
		ExploreThings bool `envconfig:"PENSIEVE_FEATURE_EXPLORE_THINGS" default:"false"`
		ExplorePeople bool `envconfig:"PENSIEVE_FEATURE_EXPLORE_PEOPLE" default:"false"`
		Sharing       bool `envconfig:"PENSIEVE_FEATURE_SHARING" default:"false"`
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
