package storage

// Disk ...
type Disk struct {
	Root string
}

func (d *Disk) Upload(fileType, filePath string) (string, error) {
	return "", nil
}

func (d *Disk) Delete(filePath string) error {
	return nil
}

func (d *Disk) Get(filePath string) (string, error) {
	return "", nil
}
