package ddownload

import (
	"time"

	"github.com/google/uuid"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type File struct {
	// Unique file identifier (UUID)
	FileID uuid.UUID

	// Associated user identifier (UUID)
	UserID *uuid.UUID

	// Status
	Status dtypes.FileStatus

	// Media URL
	MediaUrl string

	// Media title
	MediaTitle string

	// Media description
	MediaDescription *string

	// Channel ID
	ChannelID *string

	// Original file name
	FileName string

	// File extension (e.g., "mp3", "mp4")
	Ext string

	// Full file name including extension
	FullName string

	// File size (byte)
	FileSize *int64

	// Fast partial file hash (combined hash of multiple sampled blocks; not a full-file checksum)
	PartialHash *string

	// Human-readable safe full name
	SafeReadableFullName string

	// MediaInfo holds media metadata.
	MediaInfo *dtypes.MediaInfo

	// Error message
	ErrorMessage *string

	// Downloaded timestamp
	DownloadedAt *time.Time

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time

	// Related task
	DownloadTask *DownloadTask
}

func (f *File) IsYouTube() bool {
	return hostdetect.YouTube(f.MediaUrl)
}

func (src *File) Copy() *File {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)
	copy.UserID = uptr.Copy(src.UserID)
	copy.MediaDescription = uptr.Copy(src.MediaDescription)
	copy.ChannelID = uptr.Copy(src.ChannelID)
	copy.FileSize = uptr.Copy(src.FileSize)
	copy.PartialHash = uptr.Copy(src.PartialHash)
	copy.MediaInfo = src.MediaInfo.Copy()
	copy.ErrorMessage = uptr.Copy(src.ErrorMessage)
	copy.DownloadedAt = uptr.Copy(src.DownloadedAt)
	copy.DeletedAt = uptr.Copy(src.DeletedAt)
	copy.DownloadTask = src.DownloadTask.Copy()

	return copy
}
