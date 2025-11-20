package mappers

import (
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (m *Mappers) MapFileDomainToFileInfoResponse(file *ddownload.File, downloadsDir string) *dto.GetFileInfoResponse {
	var youtubeTitle = file.YoutubeTitle
	if file.YoutubeTitle == "" {
		youtubeTitle = file.YoutubeUrl
	}

	var filePath string
	if file.FullName != "" {
		filePath = filepath.Join(downloadsDir, file.FullName)
	}

	return &dto.GetFileInfoResponse{
		FileId:               file.FileId,
		Status:               file.Status,
		YoutubeUrl:           file.YoutubeUrl,
		YoutubeTitle:         youtubeTitle,
		FileName:             file.FileName,
		FileExt:              file.Ext,
		FileFullName:         file.FullName,
		FilePath:             filePath,
		FileSize:             file.FileSize,
		SafeReadableFullName: file.SafeReadableFullName,
		StatusText:           uptr.Deref(file.ErrorMessage),
		CreatedAt:            file.CreatedAt,
	}
}
