package ddownload

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// MediaSourceIndex represents a media source index
type MediaSourceIndex struct {
	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID

	// User identifier
	UserID *uuid.UUID

	// Media title
	Title string

	// Description media
	Description *string

	// Visibility access level for media (public or private)
	Visibility dtypes.MediaVisibility

	// Number of completed views
	Views uint32

	// Media source creation timestamp
	SourceCreatedAt time.Time

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time
}

// Copy creates a copy of the MediaSourceIndex object
func (src *MediaSourceIndex) Copy() *MediaSourceIndex {
	if src == nil {
		return nil
	}

	copy := new(*src)
	copy.UserID = uptr.Copy(src.UserID)
	copy.Description = uptr.Copy(src.Description)

	return copy
}

// Validate validates
func (index *MediaSourceIndex) Validate() error {
	return nil
}

func (index *MediaSourceIndex) InitFromMediaDownload(download *MediaDownload) {
	index.DownloadID = download.DownloadID
	index.UserID = download.UserID

	index.Title = download.MediaTitle
	index.Description = download.MediaDescription

	index.Visibility = download.Visibility
	index.SourceCreatedAt = download.CreatedAt
}

func (index *MediaSourceIndex) NeedsUpdateFromMediaDownload(download *MediaDownload) bool {
	return index.UserID != download.UserID ||
		index.Title != download.MediaTitle ||
		index.Description != download.MediaDescription ||
		index.Visibility != download.Visibility
}
