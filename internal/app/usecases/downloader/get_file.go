package downloader

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

func (uc *YouTubeDownloader) GetFileInfo(
	ctx context.Context,
	userID uuid.UUID,
	fileID uuid.UUID,
) (*dto.GetFileInfoResponse, error) {
	var accessByUserID *uuid.UUID
	if uc.historyMode != dtypes.HistoryModeGlobal {
		accessByUserID = &userID
	}

	resp, err := uc.findStateAndFileInfo(ctx, accessByUserID, &fileID, nil, false)
	if err != nil {
		uc.logger.Error("Failed get file info", "error", err)
		return nil, errorx.NewByErr(err, exceptionx.ERROR)
	}
	if resp == nil {
		return nil, errorx.NewByExceptionType("file not found", exceptionx.NOT_FOUND)
	}
	return resp, nil
}

func (uc *YouTubeDownloader) findStateAndFileInfo(
	ctx context.Context,
	userID *uuid.UUID,
	fileId *uuid.UUID,
	file *ddownload.File,
	checkNotFound bool,
) (*dto.GetFileInfoResponse, error) {
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

	var (
		fileResp   *ddownload.File
		dlProgress *ddownload.DownloadProgress
	)

	state, _ := uc.dlStateCache.FindByFileId(ctx, userID, id)
	if state != nil && state.File != nil {
		fileResp = uptr.Copy(state.File)
		dlProgress = uptr.Copy(state.Progress)
	}

	if fileResp == nil && file != nil {
		fileResp = uptr.Any(*file)
	}

	if fileResp == nil {
		var (
			file *ddownload.File
			err  error
		)

		if checkNotFound {
			file, err = uc.file.GetByFileId(ctx, userID, id)
			if err != nil {
				return nil, err
			}
		} else {
			file, err = uc.file.FindByFileId(ctx, userID, id)
			if err != nil {
				return nil, err
			}
		}

		if file == nil {
			return nil, nil
		}

		fileResp = uptr.Any(*file)
	}

	return uc.mappers.MapFileDomainToFileInfoResponse(fileResp, dlProgress, uc.downloadsDir), nil
}

func (uc *YouTubeDownloader) LoadHistory(
	ctx context.Context,
	userID uuid.UUID,
	before time.Time,
	limit uint64,
) ([]*dto.GetFileInfoResponse, error) {
	if uc.historyMode == dtypes.HistoryModeDisabled {
		return []*dto.GetFileInfoResponse{}, nil
	}

	var filterUserID *uuid.UUID
	if uc.historyMode == dtypes.HistoryModePerUser {
		filterUserID = &userID
	}

	return uc.getFilesInfo(ctx, filterUserID, before, limit)
}

func (uc *YouTubeDownloader) getFilesInfo(
	ctx context.Context,
	userID *uuid.UUID,
	before time.Time,
	limit uint64,
) ([]*dto.GetFileInfoResponse, error) {
	var resps []*dto.GetFileInfoResponse

	files, err := uc.file.GetBeforeTime(ctx, userID, before, limit)
	if err != nil {
		uc.logger.Warn("Failed get files", "before", before, "limit", limit, "error", err)
		return nil, err
	}

	resps = make([]*dto.GetFileInfoResponse, 0, len(files))
	for _, file := range files {
		resp, err := uc.findStateAndFileInfo(ctx, userID, nil, file, true)
		if err != nil {
			continue
		}
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *YouTubeDownloader) GetFilePath(ctx context.Context, fileId uuid.UUID) (string, error) {
	file, err := uc.file.GetByFileId(ctx, nil, fileId)
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
func (uc *YouTubeDownloader) GetDownloadFileName(
	ctx context.Context,
	userID uuid.UUID,
	fileId uuid.UUID,
) (string, string, error) {
	var accessByUserID *uuid.UUID
	if uc.historyMode != dtypes.HistoryModeGlobal {
		accessByUserID = &userID
	}

	file, err := uc.file.GetByFileId(ctx, accessByUserID, fileId)
	if err != nil {
		uc.logger.Error("Failed find file", "error", err)
		return "", "", err
	}

	return file.SafeReadableFullName, file.Ext, nil
}
