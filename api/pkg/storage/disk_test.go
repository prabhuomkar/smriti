package storage

import (
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
			Name:   "error due to cannot open file",
			Config: &Config{Provider: "disk"},
			MockFiles: func() (string, func()) {
				return "", func() {}
			},
			WantErr:     true,
			ErrContains: "error uploading file to disk as cannot open file",
		},
		{
			Name:   "error due to cannot create file",
			Config: &Config{Provider: "disk", Root: "invalid"},
			MockFiles: func() (string, func()) {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return file.Name(), func() {
					os.Remove(file.Name())
				}
			},
			WantErr:     true,
			ErrContains: "error uploading file to disk as cannot create file",
		},
		{
			Name:   "success",
			Config: &Config{Provider: "disk", Root: os.TempDir()},
			MockFiles: func() (string, func()) {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return file.Name(), func() {
					os.Remove(file.Name())
				}
			},
			WantErr:     false,
			ErrContains: "",
		},
	}
	for _, test := range tests {
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
			Name:   "error",
			Config: &Config{Provider: "disk", Root: "invalid"},
			MockFiles: func() func() {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return func() {
					os.Remove(file.Name())
				}
			},
			WantErr:     true,
			ErrContains: "error deleting file from disk",
		},
		{
			Name:   "success",
			Config: &Config{Provider: "disk", Root: os.TempDir()},
			MockFiles: func() func() {
				file, _ := os.CreateTemp(os.TempDir(), "file")
				return func() {
					os.Remove(file.Name())
				}
			},
			WantErr:     false,
			ErrContains: "",
		},
	}
	for _, test := range tests {
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
	}
}

func TestDiskGet(t *testing.T) {
	provider := Init(&Config{Provider: "disk", Root: "../storage"})
	res, err := provider.Get("fileType", "fileID")
	assert.NoError(t, err)
	assert.Equal(t, "../storage/fileType/fileID", res)
}
