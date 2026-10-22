package dto

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type ScheduleUploadRequest struct {
	Title       string
	Description string
	Visibility  dtypes.MediaVisibility
}

type ScheduleUploadResponse struct {
	DownloadID uuid.UUID
	Status     dtypes.MediaDownloadStatus
	MediaTitle string
	MediaInfo  *dtypes.MediaInfo
}
