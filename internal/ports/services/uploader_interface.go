package pservices

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type UploadedFile interface {
	// FileID returns the unique identifier of the uploaded file.
	FileID() uuid.UUID

	// OwnerID returns the unique identifier of the user who uploaded the file.
	OwnerID() uuid.UUID

	// MediaInfo returns the media information of the uploaded file.
	MediaInfo() *dtypes.MediaInfo

	// FilePath returns the full path of the uploaded file, including its extension.
	FilePath() string

	// Cleanup removes the temporary uploaded file after it has been moved to the storage.
	// An already missing file is not an error.
	Cleanup() error
}
