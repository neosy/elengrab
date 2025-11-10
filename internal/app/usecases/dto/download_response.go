package dto

import "github.com/google/uuid"

type DownloadResponse struct {
	Title string
	// format type
	Format string
	// file id
	FileId uuid.UUID
}
