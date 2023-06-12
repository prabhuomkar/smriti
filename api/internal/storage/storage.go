package storage

// Storage ...
type Storage interface {
	Upload(string, string) (string, error)
	Delete(string) error
	Get(string) (string, error)
}
