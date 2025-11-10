package ucdownloader

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (uc *YouTubeDownloader) GetFileInfo(ctx context.Context, fileId uuid.UUID) (*dto.GetFileInfoResponse, error) {
	file, err := uc.fileRep.FindByFileId(ctx, fileId)
	if err != nil {
		uc.logger.Error("Error finding record", "error", err)
		return nil, err
	}

	if file == nil {
		uc.logger.Error("Record not found", "fileId", fileId)
		return nil, err
	}

	return &dto.GetFileInfoResponse{
		Title:                file.Title,
		FileId:               file.FileId,
		Name:                 file.FileName,
		Ext:                  file.Ext,
		FullName:             file.FullName,
		Path:                 filepath.Join(uc.downloadsDir, file.FullName),
		SafeReadableFullName: file.SafeReadableFullName,
	}, nil
}

func (uc *YouTubeDownloader) GetFilePath(ctx context.Context, fileId uuid.UUID) (string, error) {
	file, err := uc.fileRep.FindByFileId(ctx, fileId)
	if err != nil {
		uc.logger.Error("Error finding record", "error", err)
		return "", err
	}

	if file == nil {
		uc.logger.Error("Record not found", "fileId", fileId)
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
	file, err := uc.fileRep.FindByFileId(ctx, fileId)
	if err != nil {
		uc.logger.Error("Error finding record", "error", err)
		return "", "", err
	}

	if file == nil {
		uc.logger.Error("Record not found", "fileId", fileId)
		return "", "", err
	}

	return file.SafeReadableFullName, file.Ext, nil
}
