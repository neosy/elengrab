package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/utils/hash"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type GetMediaDownloadInfoResponse struct {
	DownloadID    uuid.UUID
	Status        dtypes.MediaDownloadStatus
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

func (downloadInfo *GetMediaDownloadInfoResponse) IsYouTube() bool {
	return hostdetect.YouTube(downloadInfo.MediaURL)
}

func (downloadInfo *GetMediaDownloadInfoResponse) ImageMetaHash(withValues ...any) string {
	values := []any{
		downloadInfo.DownloadID,
		downloadInfo.UpdatedAt,
		downloadInfo.Status.String(),
		downloadInfo.WorkingStatus.String(),
		downloadInfo.ChannelID,
	}

	if downloadInfo.MediaInfo != nil {
		values = append(values, downloadInfo.MediaInfo.ThumbnailID)
		values = append(values, downloadInfo.MediaInfo.FrameThumbnailID)
	}

	if len(withValues) > 0 {
		values = append(values, withValues...)
	}

	return hash.MetaHashHex32(values)
}
