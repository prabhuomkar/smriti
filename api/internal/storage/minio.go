package storage

// Minio ...
type Minio struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func (m *Minio) Upload(fileType, filePath string) (string, error) {
	return "", nil
}

func (m *Minio) Delete(filePath string) error {
	return nil
}

func (m *Minio) Get(filePath string) (string, error) {
	return "", nil
}
