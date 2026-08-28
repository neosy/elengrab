package searchindex

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *SearchIndex) CreateMediaDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	index := &ddownload.MediaSourceIndex{}

	index.InitFromMediaDownload(download)

	err := uc.searchIndex.Insert(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SearchIndex) SaveMediaDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	return uc.searchIndex.Tx(ctx, func(ctx context.Context) error {
		index, err := uc.searchIndex.FindByDownload(ctx, download.DownloadID)
		if err != nil {
			return err
		}

		if index == nil {
			return uc.CreateMediaDownload(ctx, download)
		}

		if !index.NeedsUpdateFromMediaDownload(download) {
			return nil
		}

		index.InitFromMediaDownload(download)

		err = uc.searchIndex.Update(ctx, index)
		if err != nil {
			return err
		}

		return nil
	})
}
