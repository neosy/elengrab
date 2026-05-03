package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type storage struct {
	basePath string
}

func NewStorage(basePath string) (*storage, error) {
	if basePath == "" {
		return nil, fmt.Errorf("basePath is empty")
	}

	if err := validateDirPath(basePath); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &storage{basePath: basePath}, nil
}

func (s *storage) BasePath() string {
	return s.basePath
}

func (s *storage) Put(keyPath string, data []byte) error {
	path := filepath.Join(s.basePath, keyPath)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (s *storage) Get(keyPath string) ([]byte, error) {
	path := filepath.Join(s.basePath, keyPath)
	return os.ReadFile(path)
}

// Exists checks if file exists in storage.
func (s *storage) Exists(keyPath string) (bool, error) {
	path := filepath.Join(s.basePath, keyPath)

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func (s *storage) Delete(keyPath string) error {
	path := filepath.Join(s.basePath, keyPath)
	return os.Remove(path)
}

func (s *storage) Path(keyPath string) string {
	return filepath.Join(s.basePath, keyPath)
}
