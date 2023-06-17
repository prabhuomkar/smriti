package storage

import (
	"fmt"
	"io"
	"os"
)

// Disk ...
type Disk struct {
	Root string
}

func (d *Disk) Type() string {
	return ProviderDisk
}

func (d *Disk) Upload(filePath, fileType, fileID string) (string, error) {
	result := fmt.Sprintf("%s/%s/%s", d.Root, fileType, fileID)
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error uploading file to disk as cannot open file: %w", err)
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(result)
	if err != nil {
		return "", fmt.Errorf("error uploading file to disk as cannot create file: %w", err)
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("error uploading file to disk as cannot copy contents: %w", err)
	}
	defer os.Remove(filePath)
	return result, nil
}

func (d *Disk) Delete(fileType, fileID string) error {
	err := os.Remove(fmt.Sprintf("%s/%s/%s", d.Root, fileType, fileID))
	if err != nil {
		return fmt.Errorf("error deleting file from disk: %w", err)
	}
	return nil
}

func (d *Disk) Get(fileType, fileID string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", d.Root, fileType, fileID), nil
}
