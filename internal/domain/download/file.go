package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type File struct {
	// Unique file identifier (UUID)
	FileId uuid.UUID

	// Status
	Status dtypes.FileStatus

	// Youtube title
	YoutubeTitle string

	// Original file name
	FileName string

	// File extension (e.g., "mp3", "mp4")
	Ext string

	// Full file name including extension
	FullName string

	// Human-readable safe full name
	SafeReadableFullName string

	ErrorMessage *string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time

	DownloadTask *DownloadTask
}
