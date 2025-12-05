package ytdownloader

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *YouTubeDownloader) GetFileInfo(ctx context.Context, fileId uuid.UUID) (*dto.GetFileInfoResponse, error) {
	resp, err := uc.getFileInfo(ctx, nil, &fileId)
	if err != nil {
		uc.logger.Error("Failed get file info", "error", err)
		return nil, err
	}
	return resp, nil
}

func (uc *YouTubeDownloader) getFileInfo(ctx context.Context, file *ddownload.File, fileId *uuid.UUID) (*dto.GetFileInfoResponse, error) {
	var id uuid.UUID

	if file != nil {
		id = file.FileId
	} else if fileId != nil {
		id = *fileId
	}

	if id == uuid.Nil {
		uc.logger.Warn("Id for the FileId field is not defined")
		return nil, errors.New("fileId not specified")
	}

	state, _ := uc.dlState.FindByFileId(ctx, id)

	if state == nil || state.File == nil {
		f := file
		if f == nil {
			var err error
			f, err = uc.file.FindByFileId(ctx, id, true)
			if err != nil {
				uc.logger.Warn("Failed find file", "fileId", id, "error", err)
				return nil, err
			}
		}
		state = &ddownload.DownloadState{}
		state.InitFromFile(f)
	}

	return uc.mappers.MapFileDomainToFileInfoResponse(state.File, uc.downloadsDir), nil
}

func (uc *YouTubeDownloader) LoadHistory(ctx context.Context, before time.Time, limit uint64) ([]*dto.GetFileInfoResponse, error) {
	if !uc.loadHistory {
		return []*dto.GetFileInfoResponse{}, nil
	}

	return uc.getFilesInfo(ctx, before, limit)
}

func (uc *YouTubeDownloader) getFilesInfo(ctx context.Context, before time.Time, limit uint64) ([]*dto.GetFileInfoResponse, error) {
	var resps []*dto.GetFileInfoResponse

	files, err := uc.file.GetBeforeTime(ctx, before, limit)
	if err != nil {
		uc.logger.Warn("Failed get files", "before", before, "limit", limit, "error", err)
		return nil, err
	}

	resps = make([]*dto.GetFileInfoResponse, 0, len(files))
	for _, file := range files {
		resp, err := uc.getFileInfo(ctx, file, nil)
		if err != nil {
			continue
		}
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *YouTubeDownloader) GetFilePath(ctx context.Context, fileId uuid.UUID) (string, error) {
	file, err := uc.file.FindByFileId(ctx, fileId, true)
	if err != nil {
		uc.logger.Error("Failed find file", "error", err)
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
	file, err := uc.file.FindByFileId(ctx, fileId, true)
	if err != nil {
		uc.logger.Error("Failed find file", "error", err)
		return "", "", err
	}

	return file.SafeReadableFullName, file.Ext, nil
}
