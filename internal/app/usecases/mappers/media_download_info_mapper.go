package mappers

import (
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
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
		mediaTitle = download.MediaURL
	}

	workingStatus := dto.WorkingStatusNone
	if download.Status == dtypes.MediaDownloadStatusWorking {
		workingStatus = dto.WorkingStatusStartDownload
		if mappingData.Progress != nil {
			if mappingData.Progress.Percent() < 100 {
				workingStatus = dto.WorkingStatusDownloading
			} else {
				workingStatus = dto.WorkingStatusFinishDownload
			}
		}
	}

	mediaInfoText := mediaInfoText(download.MediaInfo)

	return &dto.MediaDownloadInfo{
		DownloadID: download.DownloadID,

		Status:        download.Status,
		WorkingStatus: workingStatus,

		ChannelID:   download.ChannelID,
		AvatarTitle: mappingData.AvatarTitle,

		MediaURL: download.MediaURL,

		MediaTitle:       mediaTitle,
		MediaDescription: uptr.Deref(download.MediaDescription),

		CreatedTimeAgo: humanize.TimeAgo(download.CreatedAt),

		ViewCount: mappingData.ViewCount,

		UserLastWatchPosition: mappingData.UserLastWatchPosition,
		UserWatched:           mappingData.UserWatched,

		HasSiteIcon:        mappingData.HasSiteIcon,
		ThumbnalIsPortrait: mappingData.ThumbnailIsPortrait,

		FileName:             download.FileName,
		FileExt:              download.Ext,
		FileFullName:         download.FileFullName,
		FileSize:             download.FileSize,
		SafeReadableFullName: download.SafeReadableFileFullName,
		StatusText:           uptr.Deref(download.ErrorMessage),
		MediaInfo:            download.MediaInfo,
		MediaInfoText:        mediaInfoText,
		MediaInfoTooltip:     mediaInfoTooltip(mediaInfoText),
		Progress:             mappingData.Progress,
		Visibility:           download.Visibility,

		UserID:    download.UserID,
		UserLogin: strings.ToLower(mappingData.UserLogin),

		CreatedAt: download.CreatedAt,
		UpdatedAt: download.UpdatedAt,

		MediaDownload: download,
	}
}

func mediaInfoText(mediaInfo *dtypes.MediaInfo) string {
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
