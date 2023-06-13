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

func (d *Disk) Upload(filePath, fileType, fileID string) (string, error) {
	result := fmt.Sprintf("%s/%s/%s", d.Root, fileType, fileID)
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error uploading file as cannot open file: %w", err)
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(result)
	if err != nil {
		return "", fmt.Errorf("error uploading file as cannot create file: %w", err)
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", fmt.Errorf("error uploading file as cannot copy contents: %w", err)
	}
	err = os.Remove(filePath)
	if err != nil {
		return "", fmt.Errorf("error uploading file as cannot remove file: %w", err)
	}
	return result, nil
}

func (d *Disk) Delete(filePath string) error {
	//nolint: ifshort
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	return nil
}

func (d *Disk) Get(filePath string) (string, error) {
	return filePath, nil
}
