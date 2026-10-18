package fsstorage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

// Move moves the file at fromFilePath into storage under uniqueFileName,
// creating intermediate directories as needed.
// The source file is removed on success. An existing file with the same
// storage key is overwritten.
func (s *downloadsStorage) Move(fromFilePath, uniqueFileName string) error {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	toFilePath := filepath.Join(s.BasePath(), keyPath)

	if err := os.MkdirAll(filepath.Dir(toFilePath), 0o755); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	if err := os.Rename(fromFilePath, toFilePath); err != nil {
		return fmt.Errorf("move file into storage: %w", err)
	}

	return nil
}

// Delete removes file for a given fileName.
func (s *downloadsStorage) Delete(uniqueFileName string) error {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	return s.storage.Delete(keyPath)
}

// Path returns the filesystem path where the file with uniqueFileName is (or would be) stored.
// It does not check that the file exists.
func (s *downloadsStorage) Path(uniqueFileName string) string {
	keyPath := s.buildStorageKeyPath(uniqueFileName)
	return s.storage.Path(keyPath)
}

// BasePath returns the root directory of the downloads storage.
func (s *downloadsStorage) BasePath() string {
	return s.storage.BasePath()
}

// MediaPath returns the directory where downloaded media files are stored.
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
