package mappers

import (
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapFileDomainToFileInfoResponse(
	file *ddownload.File,
	avatarTitle string,
	progress *dservices.DownloadProgress,
	hasSiteIcon bool,
) *dto.GetFileInfoResponse {
	var mediaTitle = file.MediaTitle
	if file.MediaTitle == "" {
		mediaTitle = file.MediaUrl
	}

	workingStatus := dto.WorkingStatusNone
	if file.Status == dtypes.FileStatusWorking {
		workingStatus = dto.WorkingStatusStartDownload
		if progress != nil {
			if progress.Percent() < 100 {
				workingStatus = dto.WorkingStatusDownloading
			} else {
				workingStatus = dto.WorkingStatusFinishDownload
			}
		}
	}

	mediaInfoText := mediaInfoText(file.MediaInfo)

	return &dto.GetFileInfoResponse{
		FileID:               file.FileID,
		Status:               file.Status,
		WorkingStatus:        workingStatus,
		ChannelID:            file.ChannelID,
		AvatarTitle:          avatarTitle,
		MediaUrl:             file.MediaUrl,
		MediaTitle:           mediaTitle,
		MediaDescription:     uptr.Deref(file.MediaDescription),
		HasSiteIcon:          hasSiteIcon,
		FileName:             file.FileName,
		FileExt:              file.Ext,
		FileFullName:         file.FullName,
		FileSize:             file.FileSize,
		SafeReadableFullName: file.SafeReadableFullName,
		StatusText:           uptr.Deref(file.ErrorMessage),
		MediaInfo:            file.MediaInfo,
		MediaInfoText:        mediaInfoText,
		MediaInfoTooltip:     mediaInfoTooltip(mediaInfoText),
		Progress:             progress,
		UserID:               file.UserID,
		CreatedAt:            file.CreatedAt,
		UpdatedAt:            file.UpdatedAt,
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
