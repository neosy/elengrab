package mappers

import (
	"fmt"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (m *Mappers) MapFileDomainToFileInfoResponse(file *ddownload.File, downloadsDir string) *dto.GetFileInfoResponse {
	var youtubeTitle = file.YoutubeTitle
	if file.YoutubeTitle == "" {
		youtubeTitle = file.YoutubeUrl
	}

	return &dto.GetFileInfoResponse{
		FileId:               file.FileId,
		Status:               file.Status,
		YoutubeChannelID:     file.YoutubeChannelID,
		YoutubeUrl:           file.YoutubeUrl,
		YoutubeTitle:         youtubeTitle,
		FileName:             file.FileName,
		FileExt:              file.Ext,
		FileFullName:         file.FullName,
		FileSize:             file.FileSize,
		SafeReadableFullName: file.SafeReadableFullName,
		StatusText:           uptr.Deref(file.ErrorMessage),
		MediaInfo:            file.MediaInfo,
		MediaInfoText:        mediaInfoText(file.MediaInfo),
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
			return fmt.Sprintf(
				"%v, %v, %dx%d",
				videoInfo.Codec.Title(),
				videoInfo.Resolution,
				videoInfo.Width,
				videoInfo.Height,
			)
		}
		return ""
	}

	if mediaInfo.AudioInfo != nil {
		if mediaInfo.AudioInfo.Codec != dtypes.AudioCodecNone {
			audioInfo := mediaInfo.AudioInfo
			text := fmt.Sprintf(
				"%v, %d kbps",
				audioInfo.Codec.Title(),
				audioInfo.Bitrate,
			)
			if audioInfo.SampleRate != nil {
				text += fmt.Sprintf(", %d Hz", *audioInfo.SampleRate)
			}
			return text
		}
		return ""
	}

	return ""
}
