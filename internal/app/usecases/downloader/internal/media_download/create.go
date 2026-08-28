package mediadownload

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	iconfig "github.com/neosy/elengrab/internal/config"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaDownload) Create(ctx context.Context, download *ddownload.MediaDownload, dlOptions *ddownload.DownloadOptions) error {
	if download == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if download.DownloadID == uuid.Nil {
		download.DownloadID = uuid.New()
	}

	download.Visibility = iconfig.InitialMediaVisibility()

	download.Status = dtypes.MediaDownloadStatusNew

	download.NormalizeForSave()

	err := uc.downloadRepo().Tx(ctx, func(ctx context.Context) error {
		err := uc.downloadRepo().Insert(ctx, download)
		if err != nil {
			uc.logger.Warn(
				"Failed to insert record into repository",
				"error", err,
			)
			return errorx.Errorf("failed to insert download: %w", err, exceptionx.ERROR)
		}

		download, err = uc.GetByDownloadID(ctx, download.DownloadID)
		if err != nil {
			return err
		}

		err = uc.CreateTask(ctx, download, dlOptions)
		if err != nil {
			return errorx.Errorf("failed to create task: %w", err, exceptionx.ERROR)
		}

		return nil
	})
	if err != nil {
		return err
	}

	uc.createDependencies(ctx, download)

	return nil
}

func (uc *MediaDownload) createDependencies(ctx context.Context, download *ddownload.MediaDownload) error {
	go func() {
		uc.searchIndex.CreateMediaDownload(ctx, download)
	}()
	return nil
}
