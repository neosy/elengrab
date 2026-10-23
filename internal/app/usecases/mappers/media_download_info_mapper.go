package mappers

import (
	"fmt"
	"strings"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapDownloadDomainToDownloadInfoResponse(
	download *ddownload.MediaDownload,
	mappingData *dto.MediaDownloadInfoMappingData,
) *dto.MediaDownloadInfo {
	var mediaTitle = download.MediaTitle
	if download.MediaTitle == "" {
		mediaTitle = download.SourceURL()
	}

	workingStatus := resolveWorkingStatus(download, mappingData.Progress)
	mediaInfoText := resolveMediaInfoText(download.MediaInfo)

	var userID uuid.UUID
	if download.UserID != nil {
		userID = *download.UserID
	}

	return &dto.MediaDownloadInfo{
		DownloadID: download.DownloadID,

		Status:        download.Status,
		WorkingStatus: workingStatus,

		Channel:      mappingData.Channel,
		ChannelTitle: mappingData.ChannelTitle,

		MediaURL: download.SourceURL(),

		MediaTitle:       mediaTitle,
		MediaDescription: uptr.Deref(download.MediaDescription),

		CreatedTimeAgo: humanize.TimeAgo(download.CreatedAt),

		ViewCount: mappingData.ViewCount,

		UserLastWatchPosition: mappingData.UserLastWatchPosition,
		UserWatched:           mappingData.UserWatched,

		HasSiteIcon:        mappingData.HasSiteIcon,
		ThumbnalIsPortrait: mappingData.ThumbnailIsPortrait,

		FileNameNoExt:        download.FileNameNoExt,
		FileExt:              download.Ext,
		FileName:             download.FileName,
		FileSize:             download.FileSize,
		SafeReadableFullName: download.SafeReadableFileFullName,
		StatusText:           uptr.Deref(download.ErrorMessage),
		MediaInfo:            download.MediaInfo,
		MediaInfoText:        mediaInfoText,
		MediaInfoTooltip:     mediaInfoTooltip(mediaInfoText),
		Progress:             mappingData.Progress,
		Visibility:           download.Visibility,

		UserID:    userID,
		UserType:  mappingData.UserType,
		UserLogin: strings.ToLower(mappingData.UserLogin),

		CreatedAt: download.CreatedAt,
		UpdatedAt: download.UpdatedAt,

		MediaDownload: download,
	}
}

func resolveMediaInfoText(mediaInfo *dtypes.MediaInfo) string {
	if mediaInfo == nil {
		return ""
	}

	var parts []string

	if mediaInfo.VideoInfo != nil && mediaInfo.VideoInfo.Codec != dtypes.VideoCodecNone {
		videoInfo := mediaInfo.VideoInfo
		text := fmt.Sprintf(
			"%v",
			videoInfo.Codec.Title(),
		)
		if videoInfo.Resolution != dtypes.VideoResolutionNone {
			text += fmt.Sprintf(", %v", videoInfo.Resolution)
		}
		if videoInfo.Width != 0 && videoInfo.Height != 0 {
			text += fmt.Sprintf(", %v", videoInfo.ResolutionString())
		}
		if text != "" {
			parts = append(parts, text)
		}
	}

	if mediaInfo.AudioInfo != nil && mediaInfo.AudioInfo.Codec != dtypes.AudioCodecNone {
		audioInfo := mediaInfo.AudioInfo
		text := audioInfo.Codec.Title()
		if audioInfo.Bitrate != 0 {
			text += fmt.Sprintf(", %d kbps", audioInfo.Bitrate)
		}
		if audioInfo.SampleRate != nil && *audioInfo.SampleRate != 0 {
			text += fmt.Sprintf(", %d Hz", *audioInfo.SampleRate)
		}
		if text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, "; ")
}

func mediaInfoTooltip(text string) string {
	return strings.ReplaceAll(text, "; ", "\n")
}

// resolveWorkingStatus determines the working status of a media item
// from its status, source type and the current download progress.
// Uploaded media has no download stage, so it is always reported as finished.
// Progress may be nil if no progress has been reported yet.
func resolveWorkingStatus(
	download *ddownload.MediaDownload,
	progress *dtypes.DownloadProgress,
) types.WorkingStatus {
	if download == nil || download.Status != dtypes.MediaDownloadStatusWorking {
		return types.WorkingStatusNone
	}

	if download.SourceType == dtypes.MediaDownloadSourceTypeUpload {
		return types.WorkingStatusFinishDownload
	}

	// No progress reported yet: the download has just started.
	if progress == nil {
		return types.WorkingStatusStartDownload
	}

	if progress.Percent() < 100 {
		return types.WorkingStatusDownloading
	}

	return types.WorkingStatusFinishDownload
}
