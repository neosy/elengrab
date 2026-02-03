package mappers

import (
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (m *Mappers) MapFileDomainToFileInfoResponse(
	file *ddownload.File,
	progress *ddownload.DownloadProgress,
	downloadsDir string,
	hasSiteLogo bool,
) *dto.GetFileInfoResponse {
	var youtubeTitle = file.YoutubeTitle
	if file.YoutubeTitle == "" {
		youtubeTitle = file.YoutubeUrl
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

	return &dto.GetFileInfoResponse{
		FileId:               file.FileId,
		Status:               file.Status,
		WorkingStatus:        workingStatus,
		YoutubeChannelID:     file.YoutubeChannelID,
		YoutubeUrl:           file.YoutubeUrl,
		YoutubeTitle:         youtubeTitle,
		HasSiteLogo:          hasSiteLogo,
		FileName:             file.FileName,
		FileExt:              file.Ext,
		FileFullName:         file.FullName,
		FileSize:             file.FileSize,
		SafeReadableFullName: file.SafeReadableFullName,
		StatusText:           uptr.Deref(file.ErrorMessage),
		MediaInfo:            file.MediaInfo,
		MediaInfoText:        mediaInfoText(file.MediaInfo),
		Progress:             progress,
		CreatedAt:            file.CreatedAt,
		UpdatedAt:            file.UpdatedAt,
	}
}

func mediaInfoText(mediaInfo *ddownload.MediaInfo) string {
	if mediaInfo == nil {
		return ""
	}

	if mediaInfo.VideoInfo != nil {
		if mediaInfo.VideoInfo.Codec != dtypes.VideoCodecNone {
			videoInfo := mediaInfo.VideoInfo
			text := fmt.Sprintf(
				"%v, %v",
				videoInfo.Codec.Title(),
				videoInfo.Resolution,
			)
			if videoInfo.Width != 0 && videoInfo.Height != 0 {
				text += fmt.Sprintf(", %dx%d", videoInfo.Width, videoInfo.Height)
			}
			return text
		}
		return ""
	}

	if mediaInfo.AudioInfo != nil {
		if mediaInfo.AudioInfo.Codec != dtypes.AudioCodecNone {
			audioInfo := mediaInfo.AudioInfo
			text := audioInfo.Codec.Title()
			if audioInfo.Bitrate != 0 {
				text += fmt.Sprintf(", %d kbps", audioInfo.Bitrate)
			}
			if audioInfo.SampleRate != nil && *audioInfo.SampleRate != 0 {
				text += fmt.Sprintf(", %d Hz", *audioInfo.SampleRate)
			}
			return text
		}
		return ""
	}

	return ""
}
