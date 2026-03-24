package fileuc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/pkg/errorx"
	"github.com/neosy/elengrab/pkg/errorx/exceptionx"
)

func (uc *File) Create(ctx context.Context, file *ddownload.File, dlOptions *ddownload.DownloadOptions) error {
	if file == nil {
		uc.logger.Warn("Nil pointer in function")
		return errorx.New("function parameter is a null pointer", exceptionx.ERROR)
	}

	if file.FileID == uuid.Nil {
		file.FileID = uuid.New()
	}
	file.Status = dtypes.FileStatusNew

	err := uc.fileRep.Insert(ctx, file)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.NewFromError(fmt.Errorf("failed to insert file: %w", err), exceptionx.ERROR)
	}

	err = uc.CreateTask(ctx, file, dlOptions)
	if err != nil {
		return errorx.NewFromError(fmt.Errorf("failed to create task: %w", err), exceptionx.ERROR)
	}

	uc.saveToDownloadStateCache(ctx, file.FileID)

	return nil
}
