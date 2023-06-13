package storage

import (
	"api/config"
	"fmt"
	"log"
	"os"
)

const dirPermission = 0o777

// Provider ...
type Provider interface {
	Upload(string, string, string) (string, error)
	Delete(string) error
	Get(string) (string, error)
}

func Init(cfg config.Storage) Provider { //nolint: ireturn
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
	}
	return nil
}
