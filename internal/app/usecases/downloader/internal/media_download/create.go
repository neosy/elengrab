package mediadownload

import (
	"context"
	"errors"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

const maxInsertAttempts = 3

func (uc *MediaDownload) Create(ctx context.Context, download *ddownload.MediaDownload, dlOptions *ddownload.DownloadOptions) error {
	if download == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	if download.DownloadID == uuid.Nil {
		download.DownloadID = uuid.New()
	}

	download.DownloadCode = idcodec.EncodeUUIDShortBase64URL(download.DownloadID)

	download.Status = dtypes.MediaDownloadStatusNew

	download.NormalizeForSave()

	err := uc.downloadRepo().Tx(ctx, func(ctx context.Context) error {
		var err error
		for attempt := range maxInsertAttempts {
			err = uc.downloadRepo().Insert(ctx, download)
			if err == nil || !errors.Is(err, apperrors.ErrIDConflict) {
				break
			}

			newID := download.DownloadID
			newID[15] += byte(attempt + 1)

			download.DownloadCode = idcodec.EncodeUUIDShortBase64URL(newID)
		}
		if err != nil {
			uc.downloadCacheRep.Delete(ctx, download.DownloadID)
			uc.logger.Warn(
				"Failed to insert record into repository",
				"error", err,
			)
			return errorx.Errorf("failed to insert download: %w", err, exceptionx.ERROR)
		}

		uc.downloadCacheRep.Delete(ctx, download.DownloadID)

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
