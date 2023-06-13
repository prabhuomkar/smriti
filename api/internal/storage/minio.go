package storage

import "log"

// Minio ...
type Minio struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func (m *Minio) Upload(filePath, fileType, fileID string) (string, error) {
	log.Println(filePath, fileType, fileID)
	return "", nil
}

func (m *Minio) Delete(filePath string) error {
	log.Println(filePath)
	return nil
}

func (m *Minio) Get(filePath string) (string, error) {
	log.Println(filePath)
	return "", nil
}
