package core

import (
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

func NewStorage(basePath string) (*Storage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &Storage{basePath: basePath}, nil
}

func (s *Storage) BasePath() string {
	return s.basePath
}

func (s *Storage) Put(keyPath string, data []byte) error {
	path := filepath.Join(s.basePath, keyPath)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (s *Storage) Get(keyPath string) ([]byte, error) {
	path := filepath.Join(s.basePath, keyPath)
	return os.ReadFile(path)
}

func (s *Storage) Delete(keyPath string) error {
	path := filepath.Join(s.basePath, keyPath)
	return os.Remove(path)
}

func (s *Storage) Path(keyPath string) string {
	return filepath.Join(s.basePath, keyPath)
}
