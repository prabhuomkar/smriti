package storage

import (
	"api/config"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const dirPermission = 0o777

// Provider ...
type Provider interface {
	Upload(string, string, string) (string, error)
	Delete(string, string) error
	Get(string, string) (string, error)
}

func Init(cfg config.Storage) Provider { //nolint: ireturn
	//nolint: nestif
	if cfg.Provider == "disk" {
		err := os.Mkdir(fmt.Sprintf("%s/originals", cfg.DiskRoot), dirPermission)
		if err != nil {
			log.Printf("error creating storage originals directory: %+v", err)
		}
		err = os.Mkdir(fmt.Sprintf("%s/previews", cfg.DiskRoot), dirPermission)
		if err != nil {
			log.Printf("error creating storage previews directory: %+v", err)
		}
		err = os.Mkdir(fmt.Sprintf("%s/thumbnails", cfg.DiskRoot), dirPermission)
		if err != nil {
			log.Printf("error creating storage thumbnails directory: %+v", err)
		}
		return &Disk{Root: cfg.DiskRoot}
	} else if cfg.Provider == "minio" {
		minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
			Secure: false,
		})
		if err != nil {
			log.Printf("error creating storage client: %+v", err)
		}
		return &Minio{
			Client: minioClient,
		}
	}
	return nil
}
