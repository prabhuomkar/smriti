package config

import "github.com/kelseyhightower/envconfig"

type (
	// Log ...
	Log struct {
		Level string `envconfig:"PENSIEVE_LOG_LEVEL" default:"INFO"`
	}

	// API ...
	API struct {
		Host string `envconfig:"PENSIEVE_API_HOST"`
		Port int    `envconfig:"PENSIEVE_API_PORT" default:"5001"`
	}

	// GRPC ...
	GRPC struct {
		Host string `envconfig:"PENSIEVE_GRPC_HOST"`
		Port int    `envconfig:"PENSIEVE_GRPC_PORT" default:"15001"`
	}

	// Database ...
	Database struct {
		LogLevel string `envconfig:"PENSIEVE_DATABASE_LOG_LEVEL" default:"ERROR"`
		Host     string `envconfig:"PENSIEVE_DATABASE_HOST" default:"db"`
		Port     int    `envconfig:"PENSIEVE_DATABASE_PORT" default:"5432"`
		Username string `envconfig:"PENSIEVE_DATABASE_USERNAME" default:"pensieve"`
		Password string `envconfig:"PENSIEVE_DATABASE_PASSWORD" default:"pensieve"`
		Name     string `envconfig:"PENSIEVE_DATABASE_NAME" default:"pensieve"`
	}

	// Worker ...
	Worker struct {
		Host string `envconfig:"PENSIEVE_WORKER_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"PENSIEVE_WORKER_PORT" default:"15002"`
	}

	// Auth ...
	Auth struct {
		Issuer     string `envconfig:"PENSIEVE_AUTH_ISSUER" default:"pensieve"`
		Audience   string `envconfig:"PENSIEVE_AUTH_AUDIENCE" default:"pensieve"`
		AccessTTL  int    `envconfig:"PENSIEVE_AUTH_ACCESS_TTL" default:"3600"`
		RefreshTTL int    `envconfig:"PENSIEVE_AUTH_REFRESH_TTL" default:"86400"`
		Secret     string `envconfig:"PENSIEVE_AUTH_SECRET" default:"pensieve"`
	}

	// Feature ...
	Feature struct {
		Favourites    bool `envconfig:"PENSIEVE_FEATURE_FAVOURITES" default:"true"`
		Hidden        bool `envconfig:"PENSIEVE_FEATURE_HIDDEN" default:"true"`
		Trash         bool `envconfig:"PENSIEVE_FEATURE_TRASH" default:"true"`
		Albums        bool `envconfig:"PENSIEVE_FEATURE_ALBUMS" default:"true"`
		Users         bool `envconfig:"PENSIEVE_FEATURE_USERS" default:"true"`
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
		Worker
		Auth
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
