package storage

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

const (
	ProviderDisk  = "disk"
	ProviderMinio = "minio"

	dirPermission  = 0o777
	fileFlag       = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	filePermission = 0o644
)

type ( // Provider ...
	Provider interface {
		Type() string
		Upload(filePath string, fileType string, fileID string) (string, error)
		Download(filePath string, fileType string, fileID string) error
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

//nolint:ireturn
func Init(cfg *Config) Provider {
	for _, dir := range []string{
		"originals", "previews", "thumbnails", "faces",
	} {
		err := os.Mkdir(cfg.Root+"/"+dir, dirPermission)
		if err != nil && !errors.Is(err, os.ErrExist) {
			slog.Error(fmt.Sprintf("error creating storage %s directory", dir),
				"error", err)
		}
	}

	return &Disk{Root: cfg.Root}
}
