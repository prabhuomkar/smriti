package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
		Upload(string, string, string) (string, error)
		Delete(string, string) error
		Get(string, string) (string, error)
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
			log.Printf("error creating storage client: %+v", err)
		}
		return &Minio{
			Client: minioClient,
		}
	}
	err := os.Mkdir(fmt.Sprintf("%s/originals", cfg.Root), dirPermission)
	if err != nil {
		log.Printf("error creating storage originals directory: %+v", err)
	}
	err = os.Mkdir(fmt.Sprintf("%s/previews", cfg.Root), dirPermission)
	if err != nil {
		log.Printf("error creating storage previews directory: %+v", err)
	}
	err = os.Mkdir(fmt.Sprintf("%s/thumbnails", cfg.Root), dirPermission)
	if err != nil {
		log.Printf("error creating storage thumbnails directory: %+v", err)
	}
	return &Disk{Root: cfg.Root}
}
