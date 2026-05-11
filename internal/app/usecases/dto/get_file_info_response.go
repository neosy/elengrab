package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/utils/hash"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type GetFileInfoResponse struct {
	FileID        uuid.UUID
	Status        dtypes.FileStatus
	WorkingStatus WorkingStatus

	ChannelID   *string
	AvatarTitle string

	MediaURL         string
	MediaTitle       string
	MediaDescription string

	CreatedTimeAgo string

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
	Progress             *dservices.DownloaderProgress
	UserID               *uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (fileInfo *GetFileInfoResponse) IsYouTube() bool {
	return hostdetect.YouTube(fileInfo.MediaURL)
}

func (fileInfo *GetFileInfoResponse) ImageMetaHash(withValues ...any) string {
	values := []any{
		fileInfo.FileID,
		fileInfo.UpdatedAt,
		fileInfo.Status.String(),
		fileInfo.WorkingStatus.String(),
		fileInfo.ChannelID,
	}

	if fileInfo.MediaInfo != nil {
		values = append(values, fileInfo.MediaInfo.ThumbnailID)
		values = append(values, fileInfo.MediaInfo.FrameThumbnailID)
	}

	if len(withValues) > 0 {
		values = append(values, withValues...)
	}

	return hash.MetaHashHex32(values)
}
