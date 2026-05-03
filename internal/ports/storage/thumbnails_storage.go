package pstorage

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// ThumbnailsStorage defines the interface for managing thumbnail storage.
type ThumbnailsStorage interface {
	// Path returns the full file path for the thumbnail based on the storage key and format.
	Path(storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) string
	// Put stores the thumbnail data for a given storageKey and format.
	Put(data []byte, storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) error
	// Get retrieves the thumbnail data for a given storageKey and format.
	Get(storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) ([]byte, error)
	// Delete removes the thumbnail data for a given storageKey and format.
	Delete(storageKey uuid.UUID, variant dtypes.ThumbnailVariant, format string) error
}
