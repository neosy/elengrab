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
	YoutubeChannelID     *string
	YoutubeUrl           string
	YoutubeTitle         string
	FileName             string
	FileExt              string
	FileFullName         string
	FileSize             *int
	SafeReadableFullName string
	StatusText           string
	MediaInfo            *ddownload.MediaInfo
	MediaInfoText        string
	CreatedAt            time.Time
}
