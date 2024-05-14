package models

import (
	"api/pkg/cache"
	"api/pkg/storage"
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/bluele/gcache"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

type mockMinioClient struct {
	wantErr bool
}

func (m *mockMinioClient) FPutObject(_ context.Context, _ string, _ string, _ string, _ minio.PutObjectOptions) (minio.UploadInfo, error) {
	return minio.UploadInfo{}, nil
}

func (m *mockMinioClient) FGetObject(_ context.Context, _ string, _ string, _ string, _ minio.GetObjectOptions) error {
	return nil
}

func (m *mockMinioClient) RemoveObject(_ context.Context, _ string, _ string, _ minio.RemoveObjectOptions) error {
	return nil
}

func (m *mockMinioClient) PresignedGetObject(_ context.Context, bucket string, object string, _ time.Duration, _ url.Values) (*url.URL, error) {
	if m.wantErr {
		return nil, errors.New("some error")
	}
	return url.Parse(fmt.Sprintf("https://minio/%s/%s", bucket, object))
}

func TestMediaItemsTableName(t *testing.T) {
	mediaItem := MediaItem{}
	assert.Equal(t, MediaItemsTable, mediaItem.TableName())
}

func TestMediaItemURLPluginGetMediaItemURL(t *testing.T) {
	tests := []struct {
		Name        string
		MockCache   func() cache.Provider
		MockStorage storage.Provider
		ExpectedURL string
		Args        []string
	}{
		{
			"success getting from cache",
			func() cache.Provider {
				mockCache := &cache.InMemoryCache{Connection: gcache.New(1024).LRU().Build()}
				mockCache.SetWithExpire("/originals/fileID", "cachedURL", 1*time.Minute)
				return mockCache
			},
			nil,
			"cachedURL",
			[]string{"SourceURL", "/originals/fileID"},
		},
		{
			"error getting from storage",
			func() cache.Provider { return &cache.InMemoryCache{Connection: gcache.New(1024).LRU().Build()} },
			&storage.Minio{Client: &mockMinioClient{wantErr: true}},
			"",
			[]string{"PreviewURL", "/previews/fileID"},
		},
		{
			"success getting from storage",
			func() cache.Provider { return &cache.InMemoryCache{Connection: gcache.New(1024).LRU().Build()} },
			&storage.Minio{Client: &mockMinioClient{wantErr: false}},
			"https://minio/previews/fileID",
			[]string{"PreviewURL", "/previews/fileID"},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedURL, (&MediaItemURLPlugin{
				Storage: test.MockStorage,
				Cache:   test.MockCache(),
			}).getMediaItemURL(test.Args[0], test.Args[1]))
		})
	}
}

func TestMediaItemGetFileType(t *testing.T) {
	assert.Equal(t, "originals", getFileType("SourceURL"))
	assert.Equal(t, "previews", getFileType("PreviewURL"))
	assert.Equal(t, "thumbnails", getFileType("ThumbnailURL"))
	assert.Equal(t, "unknown", getFileType("InvalidURL"))
}
