package fsstorage

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	core "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/internal"
	storagetypes "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/types"
)

type downloadsStorage struct {
	storage core.Storage

	// Options
	mediaDirName string
}

// newDownloadsStorage creates a new instance of downloadsStorage with the provided base path for storage.
func newDownloadsStorage(basePath, mediaDirName string) (*downloadsStorage, error) {
	storage, err := core.NewStorage(basePath)
	if err != nil {
		return nil, err
	}

	return &downloadsStorage{
		storage:      storage,
		mediaDirName: mediaDirName,
	}, nil
}

// Exists checks if file exists in storage.
func (s *downloadsStorage) Exists(uniqueFileName string) (bool, error) {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	return s.storage.Exists(keyPath)
}

// Move file from directory to storage
func (s *downloadsStorage) Move(fromFilePath, uniqueFileName string) error {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	toFilePath := filepath.Join(s.BasePath(), keyPath)

	if err := os.MkdirAll(filepath.Dir(toFilePath), 0755); err != nil {
		return err
	}

	return os.Rename(fromFilePath, toFilePath)
}

// Delete removes file for a given fileName.
func (s *downloadsStorage) Delete(uniqueFileName string) error {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	return s.storage.Delete(keyPath)
}

// Path returns the full file path.
func (s *downloadsStorage) Path(uniqueFileName string) string {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	return s.storage.Path(keyPath)
}

// BasePath
func (s *downloadsStorage) BasePath() string {
	return s.storage.BasePath()
}

func (s *downloadsStorage) MediaPath() string {
	return filepath.Join(s.storage.BasePath(), s.mediaDirName)
}

func (s *downloadsStorage) buildStorageKeyPath(fileName string) string {
	sum := sha256.Sum256([]byte(fileName))
	key := hex.EncodeToString(sum[:2])

	return filepath.Join(
		s.mediaDirName,
		key[:2],
		key[2:4],
		fileName,
	)
}

// Stats returns filesystem usage for the given path.
func (s *downloadsStorage) Stats() (storagetypes.StorageStats, error) {
	return s.storage.Stats()
}

// Used returns the total disk space consumed by all files under the base path.
// It recursively walks the directory and sums file sizes.
func (s *downloadsStorage) Used() (uint64, error) {
	return s.storage.Used()
}
