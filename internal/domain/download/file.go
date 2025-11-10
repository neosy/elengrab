package ddownload

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Title string

	// Unique file identifier (UUID)
	FileId uuid.UUID

	// Original file name
	FileName string

	// File extension (e.g., "mp3", "mp4")
	Ext string

	// Full file name including extension
	FullName string

	// Human-readable safe full name
	SafeReadableFullName string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}
