package mediadownload

import (
	"context"
	"errors"
	"fmt"

	"uuid"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

const maxInsertAttempts = 3

func (uc *MediaDownload) CreateDownload(
	ctx context.Context,
	download *ddownload.MediaDownload,
	dlOptions *ddownload.DownloadOptions,
	onSuccess func(),
) (*ddownload.MediaDownload, error) {
	downloadCreated, err := uc.create(
		ctx, download,
		func(ctx context.Context, download *ddownload.MediaDownload) error {
			err := uc.CreateDownloadTask(ctx, download, dlOptions)
			if err != nil {
				return errorx.Errorf("failed to create task: %w", err, exceptionx.ERROR)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	uc.dlStateCache.SaveByDownload(ctx, downloadCreated)

	if onSuccess != nil {
		onSuccess()
	}

	return downloadCreated, nil
}

func (uc *MediaDownload) CreateUpload(
	ctx context.Context,
	download *ddownload.MediaDownload,
	sourcePath string,
	onSuccess func(),
) error {
	downloadCreated, err := uc.create(
		ctx, download,
		func(ctx context.Context, download *ddownload.MediaDownload) error {
			err := uc.CreateUploadTask(ctx, download, sourcePath)
			if err != nil {
				return errorx.Errorf("failed to create task: %w", err, exceptionx.ERROR)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	uc.saveUploadStateByDownload(ctx, downloadCreated)

	if onSuccess != nil {
		onSuccess()
	}

	return nil
}

func (uc *MediaDownload) create(
	ctx context.Context,
	download *ddownload.MediaDownload,
	onCreated func(context.Context, *ddownload.MediaDownload) error,
) (*ddownload.MediaDownload, error) {
	if download == nil {
		uc.logger.Warn("Nil pointer in function")
		return nil, apperrors.ErrFuncParamNullPointer
	}

	if download.DownloadID == uuid.Nil() {
		download.DownloadID = uuid.New()
	}

	download.DownloadCode = idcodec.EncodeUUIDShortBase64URL(download.DownloadID)

	download.Status = dtypes.MediaDownloadStatusNew

	download.NormalizeForSave()

	if err := download.Validate(); err != nil {
		return nil, err
	}

	var downloadCreated *ddownload.MediaDownload

	create := func(ctx context.Context) error {
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

		err = onCreated(ctx, download)
		if err != nil {
			return err
		}

		downloadCreated, err = uc.GetByDownloadIDNoCache(ctx, download.DownloadID)
		if err != nil {
			return err
		}

		err = uc.createDependencies(ctx, downloadCreated)
		if err != nil {
			return fmt.Errorf("failed to create media download dependencies: %w", err)
		}

		return nil
	}

	err := uc.downloadRepo().Tx(ctx, create)
	if err != nil {
		return nil, err
	}

	return downloadCreated, nil
}

func (uc *MediaDownload) createDependencies(ctx context.Context, download *ddownload.MediaDownload) error {
	err := uc.searchIndex.CreateMediaDownload(ctx, download)
	if err != nil {
		return err
	}
	return nil
}
