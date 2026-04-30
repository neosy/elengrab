package fsstorage

type storage interface {
	BasePath() string
	Put(keyPath string, data []byte) error
	Get(keyPath string) ([]byte, error)
	Delete(keyPath string) error
	Path(keyPath string) string
}
