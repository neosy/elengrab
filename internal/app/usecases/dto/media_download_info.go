package dto

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/utils/hash"
	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/neosy/elengrab/internal/pkg/stringx"
)

var displayedVisibilities = map[dtypes.MediaVisibility]struct{}{
	dtypes.MediaVisibilityPrivate: {},
	dtypes.MediaVisibilityPublic:  {},
}

type MediaDownloadInfo struct {
	DownloadID    uuid.UUID
	Status        dtypes.MediaDownloadStatus
	WorkingStatus WorkingStatus

	ChannelID   *string
	AvatarTitle string

	ThumbnalIsPortrait bool

	MediaURL         string
	MediaTitle       string
	MediaDescription string

	CreatedTimeAgo string

	ViewCount uint32

	UserLastWatchPosition time.Duration
	UserWatched           bool

	HasSiteIcon          bool
	FileName             string
	FileExt              string
	FileFullName         string
	FileSize             *int64
	SafeReadableFullName string
	StatusText           string
	Progress             *dservices.DownloaderProgress

	MediaInfo        *dtypes.MediaInfo
	MediaInfoText    string
	MediaInfoTooltip string

	UserID    *uuid.UUID
	UserLogin string

	Visibility dtypes.MediaVisibility

	CreatedAt time.Time
	UpdatedAt time.Time

	HasWriteAccess bool

	MediaDownload *ddownload.MediaDownload
}

func (downloadInfo *MediaDownloadInfo) IsYouTube() bool {
	return hostdetect.YouTube(downloadInfo.MediaURL)
}

func (downloadInfo *MediaDownloadInfo) ImageMetaHash(withValues ...any) string {
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

func (info *MediaDownloadInfo) MediaDescriptionUI() string {
	description := stringx.SanitizeForMetaPreview(info.MediaDescription, 160, "...")

	if len(description) == 0 {
		description = info.MediaTitle + fmt.Sprintf(" [%s]", info.MediaInfoText)
	}

	return description
}

func (info *MediaDownloadInfo) IsReady() bool {
	return info.Status == dtypes.MediaDownloadStatusDone ||
		info.Status == dtypes.MediaDownloadStatusRefreshing
}

const (
	userWatchStartThreshold           = 8 * time.Second
	minDurationForWatchStartIndicator = 60 * time.Second
)

// WatchIndicatorEnabled reports whether the watch indicator
// is enabled for this media.
func (info *MediaDownloadInfo) WatchIndicatorEnabled() bool {
	if info.MediaInfo == nil {
		return false
	}

	// The watch start indicator is disabled for short videos.
	if info.MediaInfo.Duration() <= minDurationForWatchStartIndicator {
		return false
	}

	return true
}

// ShouldShowWatchStartIndicator reports whether the watch start indicator
// should be shown based on the user's current watch position.
func (info *MediaDownloadInfo) ShouldShowWatchStartIndicator() bool {
	if !info.WatchIndicatorEnabled() {
		return false
	}

	return info.UserLastWatchPosition >= userWatchStartThreshold
}

func (info *MediaDownloadInfo) UserWatchPercent() float64 {
	if info.MediaInfo == nil || info.MediaInfo.Duration() <= 0 {
		return 0
	}
	return math.Round(float64(info.UserLastWatchPosition)/float64(info.MediaInfo.Duration())*100*100) / 100
}

func (info *MediaDownloadInfo) UserWatchDisplayPercent() float64 {
	if !info.ShouldShowWatchStartIndicator() {
		return 0
	}

	return info.UserWatchPercent()
}

func (info *MediaDownloadInfo) IsShorts() bool {
	if info.MediaInfo == nil {
		return false
	}

	return info.MediaInfo.IsShorts()
}

func (info *MediaDownloadInfo) IsAudioOnly() bool {
	if info.MediaInfo == nil {
		return false
	}

	return info.MediaInfo.IsAudioOnly()
}

func (info *MediaDownloadInfo) ViewCountText() string {
	return fmt.Sprintf("%s views", humanize.CompactNumber(info.ViewCount))
}

func (info *MediaDownloadInfo) ShouldShowVisibility() bool {
	_, exists := displayedVisibilities[info.Visibility]
	return exists
}
