package storage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

const expiryTime = 24

// Minio ...
type Minio struct {
	Client *minio.Client
}

func (m *Minio) Upload(filePath, fileType, fileID string) (string, error) {
	contentType := "application/octet-stream"
	objectOptions := minio.PutObjectOptions{
		ContentType: contentType,
	}
	res, err := m.Client.FPutObject(context.Background(), fileType, fileID, filePath, objectOptions)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %w", err)
	}
	return res.Location, nil
}

func (m *Minio) Delete(fileType, fileID string) error {
	err := m.Client.RemoveObject(context.Background(), fileType, fileID, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("error deleting file: %w", err)
	}
	return err
}

func (m *Minio) Get(fileType, fileID string) (string, error) {
	presignedURL, err := m.Client.PresignedGetObject(context.Background(),
		fileType, fileID, expiryTime*time.Hour, url.Values{})
	if err != nil {
		log.Println("error getting file: %w", err)
	}
	return presignedURL.RequestURI(), nil
}
