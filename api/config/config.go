package config

import "github.com/kelseyhightower/envconfig"

type (
	// Log ...
	Log struct {
		Level string `envconfig:"SMRITI_LOG_LEVEL" default:"INFO"`
	}

	// API ...
	API struct {
		Host string `envconfig:"SMRITI_API_HOST"`
		Port int    `envconfig:"SMRITI_API_PORT" default:"5001"`
	}

	// GRPC ...
	GRPC struct {
		Host string `envconfig:"SMRITI_GRPC_HOST"`
		Port int    `envconfig:"SMRITI_GRPC_PORT" default:"15001"`
	}

	// Database ...
	Database struct {
		LogLevel string `envconfig:"SMRITI_DATABASE_LOG_LEVEL" default:"ERROR"`
		Host     string `envconfig:"SMRITI_DATABASE_HOST" default:"database"`
		Port     int    `envconfig:"SMRITI_DATABASE_PORT" default:"5432"`
		Username string `envconfig:"SMRITI_DATABASE_USERNAME" default:"smritiuser"`
		Password string `envconfig:"SMRITI_DATABASE_PASSWORD" default:"smritipass"`
		Name     string `envconfig:"SMRITI_DATABASE_NAME" default:"smriti"`
	}

	// Cache ...
	Cache struct {
		Type     string `envconfig:"SMRITI_CACHE_TYPE" default:"inmemory"`
		Host     string `envconfig:"SMRITI_CACHE_HOST" default:"cache"`
		Port     int    `envconfig:"SMRITI_CACHE_PORT" default:"6379"`
		Password string `envconfig:"SMRITI_CACHE_PASSWORD" default:"smritipass"`
	}

	// Worker ...
	Worker struct {
		Host string `envconfig:"SMRITI_WORKER_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"SMRITI_WORKER_PORT" default:"15002"`
	}

	// Auth ...
	Auth struct {
		Enabled    bool   `envconfig:"SMRITI_AUTH_ENABLED" default:"false"`
		Issuer     string `envconfig:"SMRITI_AUTH_ISSUER" default:"smriti"`
		Audience   string `envconfig:"SMRITI_AUTH_AUDIENCE" default:"smriti"`
		AccessTTL  int    `envconfig:"SMRITI_AUTH_ACCESS_TTL" default:"3600"`
		RefreshTTL int    `envconfig:"SMRITI_AUTH_REFRESH_TTL" default:"86400"`
		Secret     string `envconfig:"SMRITI_AUTH_SECRET" default:"smriti"`
	}

	// ML ...
	ML struct {
		Places                 bool     `envconfig:"SMRITI_ML_PLACES" default:"true"`
		Classification         bool     `envconfig:"SMRITI_ML_CLASSIFICATION" default:"false"`
		Detection              bool     `envconfig:"SMRITI_ML_DETECTION" default:"false"`
		Faces                  bool     `envconfig:"SMRITI_ML_FACES" default:"false"`
		OCR                    bool     `envconfig:"SMRITI_ML_OCR" default:"false"`
		Speech                 bool     `envconfig:"SMRITI_ML_SPEECH" default:"false"`
		PlacesProvider         string   `envconfig:"SMRITI_ML_PLACES_PROVIDER" default:"openstreetmap"`
		ClassificationDownload []string `envconfig:"SMRITI_ML_CLASSIFICATION_DOWNLOAD"`
		DetectionDownload      []string `envconfig:"SMRITI_ML_DETECTION_DOWNLOAD"`
		FacesDownload          []string `envconfig:"SMRITI_ML_FACES_DOWNLOAD"`
		OCRDownload            []string `envconfig:"SMRITI_ML_OCR_DOWNLOAD"`
		SpeechDownload         []string `envconfig:"SMRITI_ML_SPEECH_DOWNLOAD"`
	}

	// Feature ...
	Feature struct {
		Favourites bool `envconfig:"SMRITI_FEATURE_FAVOURITES" default:"true"`
		Hidden     bool `envconfig:"SMRITI_FEATURE_HIDDEN" default:"true"`
		Trash      bool `envconfig:"SMRITI_FEATURE_TRASH" default:"true"`
		Albums     bool `envconfig:"SMRITI_FEATURE_ALBUMS" default:"true"`
		Explore    bool `envconfig:"SMRITI_FEATURE_EXPLORE" default:"true"`
		Places     bool `envconfig:"SMRITI_FEATURE_PLACES" default:"true"`
		Things     bool `envconfig:"SMRITI_FEATURE_THINGS" default:"false"`
		People     bool `envconfig:"SMRITI_FEATURE_PEOPLE" default:"false"`
		Sharing    bool `envconfig:"SMRITI_FEATURE_SHARING" default:"false"`
	}

	// Admin ...
	Admin struct {
		Username string `envconfig:"SMRITI_ADMIN_USERNAME" default:"smriti"`
		Password string `envconfig:"SMRITI_ADMIN_PASSWORD" default:"smritiT3st!"`
	}

	// Storage ...
	Storage struct {
		Provider       string `envconfig:"SMRITI_STORAGE_PROVIDER" default:"disk"`
		DiskRoot       string `envconfig:"SMRITI_STORAGE_DISK_ROOT" default:"../storage"`
		MinioEndpoint  string `envconfig:"SMRITI_STORAGE_MINIO_ENDPOINT" default:"storage:9000"`
		MinioAccessKey string `envconfig:"SMRITI_STORAGE_MINIO_ACCESS_KEY" default:"smritiuser"`
		MinioSecretKey string `envconfig:"SMRITI_STORAGE_MINIO_SECRET_KEY" default:"smritipass"`
	}

	// Config ...
	Config struct {
		Log
		API
		GRPC
		Database
		Cache
		Worker
		Auth
		Feature
		ML
		Admin
		Storage
	}
)

// Init ...
func Init() (*Config, error) {
	var cfg Config
	err := envconfig.Process("smriti", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
