package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

const expiryTime = 24

type (
	// Minio ...
	Minio struct {
		Client minioClient
	}

	minioClient interface {
		FPutObject(ctx context.Context, fileType string, fileID string, filePath string, opts minio.PutObjectOptions) (minio.UploadInfo, error)
		FGetObject(ctx context.Context, fileType string, fileID string, filePath string, opts minio.GetObjectOptions) error
		RemoveObject(ctx context.Context, fileType string, fileID string, opts minio.RemoveObjectOptions) error
		PresignedGetObject(ctx context.Context, fileType string, fileID string, expiry time.Duration, vals url.Values) (*url.URL, error)
	}
)

func (m *Minio) Type() string {
	return ProviderMinio
}

func (m *Minio) Upload(filePath, fileType, fileID string) (string, error) {
	contentType := "application/octet-stream"
	objectOptions := minio.PutObjectOptions{
		ContentType: contentType,
	}
	_, err := m.Client.FPutObject(context.Background(), fileType, fileID, filePath, objectOptions)
	if err != nil {
		return "", fmt.Errorf("error uploading file to minio: %w", err)
	}
	if fileType != "previews" {
		defer os.Remove(filePath)
	}
	return fmt.Sprintf("/%s/%s", fileType, fileID), nil
}

func (m *Minio) Download(filePath, fileType, fileID string) error {
	err := m.Client.FGetObject(context.Background(), fileType, fileID, filePath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("error downloading file from minio: %w", err)
	}
	return nil
}

func (m *Minio) Delete(fileType, fileID string) error {
	err := m.Client.RemoveObject(context.Background(), fileType, fileID, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error deleting file from minio: %w", err)
	}
	return nil
}

func (m *Minio) Get(fileType, fileID string) (string, error) {
	presignedURL, err := m.Client.PresignedGetObject(context.Background(),
		fileType, fileID, expiryTime*time.Hour, url.Values{})
	if err != nil {
		return "", fmt.Errorf("error getting file from minio: %w", err)
	}
	return presignedURL.String(), nil
}
