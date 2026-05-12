package mappers

import (
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapDownloadDomainToDownloadInfoResponse(
	download *ddownload.MediaDownload,
	avatarTitle string,
	progress *dservices.DownloaderProgress,
	hasSiteIcon bool,
) *dto.GetMediaDownloadInfoResponse {
	var mediaTitle = download.MediaTitle
	if download.MediaTitle == "" {
		mediaTitle = download.MediaURL
	}

	workingStatus := dto.WorkingStatusNone
	if download.Status == dtypes.MediaDownloadStatusWorking {
		workingStatus = dto.WorkingStatusStartDownload
		if progress != nil {
			if progress.Percent() < 100 {
				workingStatus = dto.WorkingStatusDownloading
			} else {
				workingStatus = dto.WorkingStatusFinishDownload
			}
		}
	}

	mediaInfoText := mediaInfoText(download.MediaInfo)

	return &dto.GetMediaDownloadInfoResponse{
		DownloadID:           download.DownloadID,
		Status:               download.Status,
		WorkingStatus:        workingStatus,
		ChannelID:            download.ChannelID,
		AvatarTitle:          avatarTitle,
		MediaURL:             download.MediaURL,
		MediaTitle:           mediaTitle,
		MediaDescription:     uptr.Deref(download.MediaDescription),
		CreatedTimeAgo:       humanize.TimeAgo(download.CreatedAt),
		HasSiteIcon:          hasSiteIcon,
		FileName:             download.FileName,
		FileExt:              download.Ext,
		FileFullName:         download.FileFullName,
		FileSize:             download.FileSize,
		SafeReadableFullName: download.SafeReadableFullName,
		StatusText:           uptr.Deref(download.ErrorMessage),
		MediaInfo:            download.MediaInfo,
		MediaInfoText:        mediaInfoText,
		MediaInfoTooltip:     mediaInfoTooltip(mediaInfoText),
		Progress:             progress,
		UserID:               download.UserID,
		CreatedAt:            download.CreatedAt,
		UpdatedAt:            download.UpdatedAt,
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
