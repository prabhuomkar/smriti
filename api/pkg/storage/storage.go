package storage

import (
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/slog"
)

const (
	ProviderDisk  = "disk"
	ProviderMinio = "minio"

	dirPermission = 0o777
)

type (
	// Provider ...
	Provider interface {
		Type() string
		Upload(filePath string, fileType string, fileID string) (string, error)
		Delete(fileType string, fileID string) error
		Get(fileType string, fileID string) (string, error)
	}

	// Config ...
	Config struct {
		Provider  string
		Root      string
		Endpoint  string
		AccessKey string
		SecretKey string
	}
)

func Init(cfg *Config) Provider { //nolint: ireturn
	if cfg.Provider == ProviderMinio {
		minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false,
		})
		if err != nil {
			slog.Error("error creating storage client", slog.Any("error", err))
		}
		return &Minio{
			Client: minioClient,
		}
	}
	err := os.Mkdir(fmt.Sprintf("%s/originals", cfg.Root), dirPermission)
	if err != nil {
		slog.Error("error creating storage originals directory", slog.Any("error", err))
	}
	err = os.Mkdir(fmt.Sprintf("%s/previews", cfg.Root), dirPermission)
	if err != nil {
		slog.Error("error creating storage previews directory", slog.Any("error", err))
	}
	err = os.Mkdir(fmt.Sprintf("%s/thumbnails", cfg.Root), dirPermission)
	if err != nil {
		slog.Error("error creating storage thumbnails directory", slog.Any("error", err))
	}
	return &Disk{Root: cfg.Root}
}
