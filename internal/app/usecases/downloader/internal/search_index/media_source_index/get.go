package sourceindex

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *MediaSourceIndex) FindByDownload(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaSourceIndex, error) {
	download, err := uc.indexRepo().FindByDownloadID(ctx, downloadID)
	if err != nil {
		uc.logger.Warn("Failed to find record", "error", err)
		return nil, err
	}

	return download, err
}

func (uc *MediaSourceIndex) GetByDownload(
	ctx context.Context,
	downloadID uuid.UUID,
) (*ddownload.MediaSourceIndex, error) {
	download, err := uc.FindByDownload(ctx, downloadID)
	if err != nil {
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	if download == nil {
		uc.logger.Warn("MediaSourceIndex not found", "downloadID", downloadID)
		return nil, errorx.New("mediaSourceIndex not found", exceptionx.NOT_FOUND)
	}

	return download, nil
}
