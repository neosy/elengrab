package fileuc

import (
	"context"
	"time"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *File) FindByFileID(
	ctx context.Context,
	userID *uuid.UUID,
	fileId uuid.UUID,
) (*ddownload.File, error) {
	fileRep := uc.fileRep
	if userID != nil {
		fileRep = uc.fileRep.WithUser(*userID)
	}

	file, err := fileRep.FindByFileID(ctx, fileId)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return nil, err
	}

	return file, err
}

// GetByFileID
// File MUST exist — otherwise NOT_FOUND
func (uc *File) GetByFileID(
	ctx context.Context,
	userID *uuid.UUID,
	fileID uuid.UUID,
) (*ddownload.File, error) {
	file, err := uc.FindByFileID(ctx, userID, fileID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if file == nil {
		uc.logger.Warn("File not found", "fileID", fileID)
		return nil, errorx.New("file not found", exceptionx.NOT_FOUND)
	}

	return file, nil
}

func (uc *File) GetAll(ctx context.Context, includeDeleted bool) ([]*ddownload.File, error) {
	file, err := uc.fileRep.GetAll(ctx, includeDeleted)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetAllFullNames(ctx context.Context, includeDeleted bool) ([]string, error) {
	names, err := uc.fileRep.GetAllFullNames(ctx, includeDeleted)
	if err != nil {
		uc.logger.Warn("Failed to get fullNames", "error", err)
		return nil, err
	}

	return names, nil
}

func (uc *File) GetBeforeTime(ctx context.Context, before time.Time, limit uint64, filters map[string]any) ([]*ddownload.File, error) {
	fileRep := uc.fileRep
	if filters != nil {
		fileRep = uc.fileRep.WithFilters(filters)
	}

	file, err := fileRep.GetBeforeTime(ctx, before, limit)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error) {
	file, err := uc.fileRep.GetByStatus(ctx, status)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetByPartialHash(ctx context.Context, criteria ddownload.DuplicateHashRow) ([]*ddownload.File, error) {
	fileRep := uc.fileRep
	if criteria.UserID != nil {
		fileRep = uc.fileRep.WithUser(*criteria.UserID)
	}
	file, err := fileRep.GetByPartialHash(ctx, criteria.Hash)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	return file, err
}

func (uc *File) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error) {
	var files []*ddownload.File

	gFiles, err := uc.fileRep.GetWithoutPartialHash(ctx)
	if err != nil {
		uc.logger.Warn("Failed to get files", "error", err)
		return nil, err
	}

	if len(gFiles) > 0 {
		files = make([]*ddownload.File, 0, len(gFiles))
		for _, file := range gFiles {
			if file.FullName == "" {
				continue
			}
			files = append(files, file)
		}
	}

	return files[:len(files):len(files)], nil
}

func (uc *File) GetDuplicateHashes(ctx context.Context, scope dtypes.UniquenessScope) ([]ddownload.DuplicateHashRow, error) {
	rows, err := uc.fileRep.GetDuplicateHashes(ctx, scope)
	if err != nil {
		uc.logger.Warn("Failed to get dublicate hashes", "error", err)
		return nil, err
	}

	return rows, nil
}

func (uc *File) GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.File, error) {
	files, err := uc.fileRep.GetDeleted(ctx, from, to)
	if err != nil {
		uc.logger.Warn("Failed to get deleted", "fromDate", from, "toDate", to, "error", err)
		return nil, err
	}

	return files, nil
}
