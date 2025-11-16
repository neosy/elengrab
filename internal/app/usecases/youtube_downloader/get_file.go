package ytdownloader

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) findByFileId(ctx context.Context, fileId uuid.UUID, checkNotFound bool) (*ddownload.File, error) {
	return uc.file.FindByFileId(ctx, fileId, checkNotFound)
}

func (uc *YouTubeDownloader) GetFileInfo(ctx context.Context, fileId uuid.UUID) (*dto.GetFileInfoResponse, error) {
	state, _ := uc.dlState.FindByFileId(ctx, fileId)

	if state == nil || state.File == nil {
		file, err := uc.findByFileId(ctx, fileId, true)
		if err != nil {
			return nil, err
		}
		state = &ddownload.DownloadState{}
		state.InitFromFile(file)
	}

	file := state.File

	var youtubeTitle = file.YoutubeTitle
	if file.YoutubeTitle == "" {
		youtubeTitle = file.YoutubeUrl
	}

	var filePath string
	if file.FullName != "" {
		filePath = filepath.Join(uc.downloadsDir, file.FullName)
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
	}, nil
}

func (uc *YouTubeDownloader) GetFilePath(ctx context.Context, fileId uuid.UUID) (string, error) {
	file, err := uc.findByFileId(ctx, fileId, true)
	if err != nil {
		return "", err
	}

	return filepath.Join(uc.downloadsDir, file.FullName), nil
}

// GetDownloadFileName retrieves the display file name and extension
// for the given file ID.
//
// Returns:
//
//	filename - the human-readable name of the file
//	ext      - the file extension (without dot)
//	err      - an error if the record is not found or a query fails
func (uc *YouTubeDownloader) GetDownloadFileName(ctx context.Context, fileId uuid.UUID) (string, string, error) {
	file, err := uc.findByFileId(ctx, fileId, true)
	if err != nil {
		return "", "", err
	}

	return file.SafeReadableFullName, file.Ext, nil
}
