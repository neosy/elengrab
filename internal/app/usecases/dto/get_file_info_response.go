package dto

import (
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type GetFileInfoResponse struct {
	FileId               uuid.UUID
	Status               dtypes.FileStatus
	WorkingStatus        WorkingStatus
	YoutubeChannelID     *string
	AvatarTitle          string
	MediaUrl             string
	MediaTitle           string
	HasSiteLogo          bool
	FileName             string
	FileExt              string
	FileFullName         string
	FileSize             *int64
	SafeReadableFullName string
	StatusText           string
	MediaInfo            *ddownload.MediaInfo
	MediaInfoText        string
	Progress             *ddownload.DownloadProgress
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
