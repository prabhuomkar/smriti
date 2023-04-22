package config

import "github.com/kelseyhightower/envconfig"

type (
	// Log ...
	Log struct {
		Level string `envconfig:"CAROUSEL_LOG_LEVEL" default:"INFO"`
	}

	// API ...
	API struct {
		Host string `envconfig:"CAROUSEL_API_HOST"`
		Port int    `envconfig:"CAROUSEL_API_PORT" default:"5001"`
	}

	// GRPC ...
	GRPC struct {
		Host string `envconfig:"CAROUSEL_GRPC_HOST"`
		Port int    `envconfig:"CAROUSEL_GRPC_PORT" default:"15001"`
	}

	// Database ...
	Database struct {
		LogLevel string `envconfig:"CAROUSEL_DATABASE_LOG_LEVEL" default:"ERROR"`
		Host     string `envconfig:"CAROUSEL_DATABASE_HOST" default:"database"`
		Port     int    `envconfig:"CAROUSEL_DATABASE_PORT" default:"5432"`
		Username string `envconfig:"CAROUSEL_DATABASE_USERNAME" default:"carousel"`
		Password string `envconfig:"CAROUSEL_DATABASE_PASSWORD" default:"carousel"`
		Name     string `envconfig:"CAROUSEL_DATABASE_NAME" default:"carousel"`
	}

	// Worker ...
	Worker struct {
		Host string `envconfig:"CAROUSEL_WORKER_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"CAROUSEL_WORKER_PORT" default:"15002"`
	}

	// Auth ...
	Auth struct {
		Enabled    bool   `envconfig:"CAROUSEL_AUTH_ENABLED" default:"false"`
		Issuer     string `envconfig:"CAROUSEL_AUTH_ISSUER" default:"carousel"`
		Audience   string `envconfig:"CAROUSEL_AUTH_AUDIENCE" default:"carousel"`
		AccessTTL  int    `envconfig:"CAROUSEL_AUTH_ACCESS_TTL" default:"3600"`
		RefreshTTL int    `envconfig:"CAROUSEL_AUTH_REFRESH_TTL" default:"86400"`
		Secret     string `envconfig:"CAROUSEL_AUTH_SECRET" default:"carousel"`
	}

	// ML ...
	ML struct {
		Places                 bool     `envconfig:"CAROUSEL_ML_PLACES" default:"true"`
		Classification         bool     `envconfig:"CAROUSEL_ML_CLASSIFICATION" default:"false"`
		Detection              bool     `envconfig:"CAROUSEL_ML_DETECTION" default:"false"`
		Faces                  bool     `envconfig:"CAROUSEL_ML_FACES" default:"false"`
		OCR                    bool     `envconfig:"CAROUSEL_ML_OCR" default:"false"`
		Speech                 bool     `envconfig:"CAROUSEL_ML_SPEECH" default:"false"`
		PlacesSource           string   `envconfig:"CAROUSEL_ML_PLACES_SOURCE" default:"openstreetmap"`
		ClassificationDownload []string `envconfig:"CAROUSEL_ML_CLASSIFICATION_DOWNLOAD"`
		DetectionDownload      []string `envconfig:"CAROUSEL_ML_DETECTION_DOWNLOAD"`
		FacesDownload          []string `envconfig:"CAROUSEL_ML_FACES_DOWNLOAD"`
		OCRDownload            []string `envconfig:"CAROUSEL_ML_OCR_DOWNLOAD"`
		SpeechDownload         []string `envconfig:"CAROUSEL_ML_SPEECH_DOWNLOAD"`
	}

	// Feature ...
	Feature struct {
		Favourites bool `envconfig:"CAROUSEL_FEATURE_FAVOURITES" default:"true"`
		Hidden     bool `envconfig:"CAROUSEL_FEATURE_HIDDEN" default:"true"`
		Trash      bool `envconfig:"CAROUSEL_FEATURE_TRASH" default:"true"`
		Albums     bool `envconfig:"CAROUSEL_FEATURE_ALBUMS" default:"true"`
		Explore    bool `envconfig:"CAROUSEL_FEATURE_EXPLORE" default:"true"`
		Places     bool `envconfig:"CAROUSEL_FEATURE_PLACES" default:"true"`
		Things     bool `envconfig:"CAROUSEL_FEATURE_THINGS" default:"false"`
		People     bool `envconfig:"CAROUSEL_FEATURE_PEOPLE" default:"false"`
		Sharing    bool `envconfig:"CAROUSEL_FEATURE_SHARING" default:"false"`
	}

	// Admin ...
	Admin struct {
		Username string `envconfig:"CAROUSEL_ADMIN_USERNAME" default:"carousel"`
		Password string `envconfig:"CAROUSEL_ADMIN_PASSWORD" default:"carouselT3st!"`
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
		ML
		Admin
		StorageDiskRoot string `envconfig:"CAROUSEL_STORAGE_DISK_ROOT" default:"../storage"`
	}
)

// Init ...
func Init() (*Config, error) {
	var cfg Config
	err := envconfig.Process("carousel", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
