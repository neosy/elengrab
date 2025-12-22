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

	var mediaInfoText string
	if file.MediaInfo != nil {
		if file.MediaInfo.VideoCodec != dtypes.VideoCodecNone {
			mediaInfoText = fmt.Sprintf(
				"%v, %v, %dx%d",
				file.MediaInfo.VideoCodec.Title(),
				file.MediaInfo.Resolution,
				file.MediaInfo.Width,
				file.MediaInfo.Height,
			)
		}
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
		MediaInfoText:        mediaInfoText,
		CreatedAt:            file.CreatedAt,
	}
}
