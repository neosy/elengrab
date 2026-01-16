package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type File struct {
	// Unique file identifier (UUID)
	FileId uuid.UUID

	// Associated user identifier (UUID)
	UserID *uuid.UUID

	// Status
	Status dtypes.FileStatus

	// Youtube URL
	YoutubeUrl string

	// Youtube title
	YoutubeTitle string

	// Youtube Channel ID
	YoutubeChannelID *string

	// Original file name
	FileName string

	// File extension (e.g., "mp3", "mp4")
	Ext string

	// Full file name including extension
	FullName string

	// File size (byte)
	FileSize *int

	// Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
	PartialHash *string

	// Human-readable safe full name
	SafeReadableFullName string

	// MediaInfo holds media metadata.
	MediaInfo *MediaInfo

	// Error message
	ErrorMessage *string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time

	// Related task
	DownloadTask *DownloadTask
}
