package core

type Storage interface {
	Get(keyPath string) ([]byte, error)
	// Exists checks if file exists in storage.
	Exists(keyPath string) (bool, error)

	Put(keyPath string, data []byte) error
	Delete(keyPath string) error

	BasePath() string
	Path(keyPath string) string

	// Stats returns filesystem usage for the given path.
	Stats() (StorageStats, error)
}
