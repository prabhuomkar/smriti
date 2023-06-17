package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

type mockMinioClient struct {
	wantErr bool
}

func (m *mockMinioClient) FPutObject(_ context.Context, _ string, _ string, _ string, _ minio.PutObjectOptions) (minio.UploadInfo, error) {
	if m.wantErr {
		return minio.UploadInfo{}, errors.New("some error")
	}
	return minio.UploadInfo{}, nil
}

func (m *mockMinioClient) RemoveObject(_ context.Context, _ string, _ string, _ minio.RemoveObjectOptions) error {
	if m.wantErr {
		return errors.New("some error")
	}
	return nil
}

func (m *mockMinioClient) PresignedGetObject(_ context.Context, bucket string, object string, _ time.Duration, _ url.Values) (*url.URL, error) {
	if m.wantErr {
		return nil, errors.New("some error")
	}
	return url.Parse(fmt.Sprintf("https://minio/%s/%s", bucket, object))
}

func TestMinioType(t *testing.T) {
	provider := Init(&Config{Provider: "minio"})
	assert.Equal(t, ProviderMinio, provider.Type())
}

func TestMinioUpload(t *testing.T) {
	tests := []struct {
		Name        string
		WantErr     bool
		ErrContains string
	}{
		{
			"error",
			true,
			"error uploading file to minio",
		},
		{
			"success",
			false,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			provider := &Minio{&mockMinioClient{test.WantErr}}
			res, err := provider.Upload("filePath", "originals", "fileID")
			if test.WantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ErrContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "/originals/fileID", res)
			}
		})
	}
}

func TestMinioDelete(t *testing.T) {
	tests := []struct {
		Name        string
		WantErr     bool
		ErrContains string
	}{
		{
			"error",
			true,
			"error deleting file from minio",
		},
		{
			"success",
			false,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			provider := &Minio{&mockMinioClient{test.WantErr}}
			err := provider.Delete("originals", "fileID")
			if test.WantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ErrContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMinioGet(t *testing.T) {
	tests := []struct {
		Name        string
		WantErr     bool
		ErrContains string
	}{
		{
			"error",
			true,
			"error getting file from minio",
		},
		{
			"success",
			false,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			provider := &Minio{&mockMinioClient{test.WantErr}}
			res, err := provider.Get("originals", "fileID")
			if test.WantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ErrContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "https://minio/originals/fileID", res)
			}
		})
	}
}
