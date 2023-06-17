package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskType(t *testing.T) {
	provider := Init(&Config{Provider: "disk", Root: os.TempDir()})
	assert.Equal(t, ProviderDisk, provider.Type())
}

func TestDiskUpload(t *testing.T) {
	tests := []struct {
		Name        string
		Config      *Config
		MockFiles   func() (string, func())
		WantErr     bool
		ErrContains string
	}{
		{
			"error due to cannot open file",
			&Config{Provider: "disk", Root: os.TempDir()},
			func() (string, func()) {
				return "", func() {}
			},
			true,
			"error uploading file to disk as cannot open file",
		},
		{
			"error due to cannot create file",
			&Config{Provider: "disk", Root: "invalid"},
			func() (string, func()) {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return file.Name(), func() {
					os.Remove(file.Name())
				}
			},
			true,
			"error uploading file to disk as cannot create file",
		},
		{
			"success",
			&Config{Provider: "disk", Root: os.TempDir()},
			func() (string, func()) {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return file.Name(), func() {
					os.Remove(file.Name())
				}
			},
			false,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			provider := Init(test.Config)
			filePath, clear := test.MockFiles()
			defer clear()
			res, err := provider.Upload(filePath, "originals", "fileID")
			if test.WantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ErrContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "/tmp/originals/fileID", res)
			}
		})
	}
}

func TestDiskDelete(t *testing.T) {
	tests := []struct {
		Name        string
		Config      *Config
		MockFiles   func() func()
		WantErr     bool
		ErrContains string
	}{
		{
			"error",
			&Config{Provider: "disk", Root: "invalid"},
			func() func() {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return func() {
					os.Remove(file.Name())
				}
			},
			true,
			"error deleting file from disk",
		},
		{
			"success",
			&Config{Provider: "disk", Root: os.TempDir()},
			func() func() {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return func() {
					os.Remove(file.Name())
				}
			},
			false,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			provider := Init(test.Config)
			clear := test.MockFiles()
			defer clear()
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

func TestDiskGet(t *testing.T) {
	tmpRoot := os.TempDir()
	provider := Init(&Config{Provider: "disk", Root: tmpRoot})
	res, err := provider.Get("fileType", "fileID")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s/fileType/fileID", tmpRoot), res)
}
