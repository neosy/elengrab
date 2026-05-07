package dto

import (
	"time"

	"github.com/google/uuid"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type GetFileInfoResponse struct {
	FileID               uuid.UUID
	Status               dtypes.FileStatus
	WorkingStatus        WorkingStatus
	ChannelID            *string
	AvatarTitle          string
	MediaUrl             string
	MediaTitle           string
	MediaDescription     string
	HasSiteIcon          bool
	FileName             string
	FileExt              string
	FileFullName         string
	FileSize             *int64
	SafeReadableFullName string
	StatusText           string
	MediaInfo            *dtypes.MediaInfo
	MediaInfoText        string
	MediaInfoTooltip     string
	Progress             *dservices.DownloadProgress
	UserID               *uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (r *GetFileInfoResponse) IsYouTube() bool {
	return hostdetect.YouTube(r.MediaUrl)
}
