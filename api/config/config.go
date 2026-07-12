package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ( // Log ...
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
		LogLevel string        `envconfig:"SMRITI_DATABASE_LOG_LEVEL" default:"ERROR"`
		Host     string        `envconfig:"SMRITI_DATABASE_HOST"      default:"database"`
		Port     int           `envconfig:"SMRITI_DATABASE_PORT"      default:"5432"`
		Username string        `envconfig:"SMRITI_DATABASE_USERNAME"  default:"smritiuser"`
		Password string        `envconfig:"SMRITI_DATABASE_PASSWORD"  default:"smritipass"`
		Name     string        `envconfig:"SMRITI_DATABASE_NAME"      default:"smriti"`
		Timeout  time.Duration `envconfig:"SMRITI_DATABASE_TIMEOUT"   default:"10s"`
	}

	// Cache ...
	Cache struct {
		Type     string `envconfig:"SMRITI_CACHE_TYPE"     default:"inmemory"`
		Host     string `envconfig:"SMRITI_CACHE_HOST"     default:"cache"`
		Port     int    `envconfig:"SMRITI_CACHE_PORT"     default:"6379"`
		Password string `envconfig:"SMRITI_CACHE_PASSWORD" default:"smritipass"`
	}

	// Worker ...
	Worker struct {
		Host string `envconfig:"SMRITI_WORKER_HOST" default:"127.0.0.1"`
		Port int    `envconfig:"SMRITI_WORKER_PORT" default:"15002"`
	}

	// Auth ...
	Auth struct {
		Enabled    bool   `envconfig:"SMRITI_AUTH_ENABLED"     default:"false"`
		Issuer     string `envconfig:"SMRITI_AUTH_ISSUER"      default:"smriti"`
		Audience   string `envconfig:"SMRITI_AUTH_AUDIENCE"    default:"smriti"`
		AccessTTL  int    `envconfig:"SMRITI_AUTH_ACCESS_TTL"  default:"3600"`
		RefreshTTL int    `envconfig:"SMRITI_AUTH_REFRESH_TTL" default:"86400"`
		Secret     string `envconfig:"SMRITI_AUTH_SECRET"      default:"smriti"`
	}

	// ML ...
	ML struct {
		Places                 bool   `envconfig:"SMRITI_ML_PLACES"                   default:"true"`
		Search                 bool   `envconfig:"SMRITI_ML_SEARCH"                   default:"true"`
		Faces                  bool   `envconfig:"SMRITI_ML_FACES"                    default:"true"`
		PlacesProvider         string `envconfig:"SMRITI_ML_PLACES_PROVIDER"          default:"openstreetmap"`
		SearchProvider         string `envconfig:"SMRITI_ML_SEARCH_PROVIDER"          default:"pytorch"`
		SearchParams           string `envconfig:"SMRITI_ML_SEARCH_PARAMS"            default:"{\"tokenizer_dir\":\"search_tokenizer\",\"processor_dir\":\"search_processor\",\"text_file\":\"search_text_v240624.pt\",\"vision_file\":\"search_vision_v240624.pt\"}"` //nolint:lll
		FacesProvider          string `envconfig:"SMRITI_ML_FACES_PROVIDER"           default:"onnx"`
		FacesParams            string `envconfig:"SMRITI_ML_FACES_PARAMS"             default:"{\"detection_threshold\":0.8,\"detection_model\":\"faces_det/scrfd_2.5g.onnx\",\"recognition_model\":\"faces_rec/webface_r50.onnx\",\"clustering_lib\":\"faiss\",\"clustering_cron\":\"* * * * *\"}"` //nolint:lll
		PreviewThumbnailParams string `envconfig:"SMRITI_ML_PREVIEW_THUMBNAIL_PARAMS" default:"{\"image_quality\":50,\"thumbnail_size\":256,\"placeholder_size\":2}"`                                                                                                                                //nolint:lll
	}

	// Feature ...
	Feature struct {
		Favourites bool `envconfig:"SMRITI_FEATURE_FAVOURITES" default:"true"`
		Hidden     bool `envconfig:"SMRITI_FEATURE_HIDDEN"     default:"true"`
		Trash      bool `envconfig:"SMRITI_FEATURE_TRASH"      default:"true"`
		Albums     bool `envconfig:"SMRITI_FEATURE_ALBUMS"     default:"true"`
		Explore    bool `envconfig:"SMRITI_FEATURE_EXPLORE"    default:"true"`
		Places     bool `envconfig:"SMRITI_FEATURE_PLACES"     default:"true"`
		People     bool `envconfig:"SMRITI_FEATURE_PEOPLE"     default:"true"`
		Sharing    bool `envconfig:"SMRITI_FEATURE_SHARING"    default:"true"`
		Jobs       bool `envconfig:"SMRITI_FEATURE_JOBS"       default:"true"`
	}

	// Admin ...
	Admin struct {
		Username string `envconfig:"SMRITI_ADMIN_USERNAME" default:"smriti"`
		Password string `envconfig:"SMRITI_ADMIN_PASSWORD" default:"smritiT3st!"`
	}

	// Storage ...
	Storage struct {
		Provider  string `envconfig:"SMRITI_STORAGE_PROVIDER"   default:"disk"`
		DiskRoot  string `envconfig:"SMRITI_STORAGE_DISK_ROOT"  default:"../storage"`
		Endpoint  string `envconfig:"SMRITI_STORAGE_ENDPOINT"   default:"storage:9000"`
		AccessKey string `envconfig:"SMRITI_STORAGE_ACCESS_KEY" default:"smritiuser"`
		SecretKey string `envconfig:"SMRITI_STORAGE_SECRET_KEY" default:"smritipass"`
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
