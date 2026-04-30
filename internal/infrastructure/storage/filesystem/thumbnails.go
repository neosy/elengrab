package fsstorage

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	core "github.com/neosy/elengrab/internal/infrastructure/storage/filesystem/internal"
)

type ThumbnailsStorage struct {
	storage storage
}

// NewThumbnailsStorage creates a new instance of ThumbnailsStorage with the provided base path for storage.
func NewThumbnailsStorage(basePath string) (*ThumbnailsStorage, error) {
	storage, err := core.NewStorage(basePath)
	if err != nil {
		return nil, err
	}

	return &ThumbnailsStorage{
		storage: storage,
	}, nil
}

// Put stores the thumbnail data for a given storageKey and format.
func (s *ThumbnailsStorage) Put(data []byte, storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) error {
	key := s.BuildStorageKeyPath(storageKey, variant, format)
	return s.storage.Put(key, data)
}

// Get retrieves the thumbnail data for a given storageKey and format.
func (s *ThumbnailsStorage) Get(storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) ([]byte, error) {
	key := s.BuildStorageKeyPath(storageKey, variant, format)
	return s.storage.Get(key)
}

// Delete removes the thumbnail data for a given storageKey and format.
func (s *ThumbnailsStorage) Delete(storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) error {
	key := s.BuildStorageKeyPath(storageKey, variant, format)
	return s.storage.Delete(key)
}

// Path returns the full file path for the thumbnail based on the storage key and format.
func (s *ThumbnailsStorage) Path(storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) string {
	key := s.BuildStorageKeyPath(storageKey, variant, format)
	return s.storage.Path(key)
}

func (s *ThumbnailsStorage) BuildStorageKeyPath(
	storageKey uuid.UUID,
	variant dtypes.ThumbnailVariant,
	format string,
) string {
	key := strings.ReplaceAll(storageKey.String(), "-", "")
	return filepath.Join(variant.String(), core.BuildStorageKeyPath(key, format))
}
