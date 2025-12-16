package fileuc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *File) Create(ctx context.Context, file *ddownload.File, dlOptions *ddownload.DownloadOptions) error {
	if file == nil {
		uc.logger.Warn("Nil pointer in function")
		return errors.New("function parameter is a null pointer")
	}

	if file.FileId == uuid.Nil {
		file.FileId = uuid.New()
	}
	file.Status = dtypes.FileStatusNew

	err := uc.fileRep.Insert(ctx, file)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return err
	}

	err = uc.CreateTask(ctx, file, dlOptions)
	if err != nil {
		return err
	}

	uc.saveToDownloadStateCache(ctx, file.FileId)

	return nil
}
