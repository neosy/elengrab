package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *downloader) PatchMediaDownload(
	ctx context.Context,
	authCtx dauth.AuthContext,
	req dto.PatchMediaDownloadRequest,
) error {
	err := uc.validateWriteOperation(authCtx)
	if err != nil {
		return err
	}

	req.Normalize()

	if err := req.Validate(); err != nil {
		return err
	}

	download, err := uc.download.GetByDownloadID(ctx, req.DownloadID)
	if err != nil {
		return err
	}

	err = uc.validateDownloadEditAccess(authCtx, download)
	if err != nil {
		return err
	}

	err = uc.download.UpdateFields(ctx, req)
	if err != nil {
		return err
	}

	downloadChanged := &dto.MediaDownloadChanged{
		DownloadID: download.DownloadID,
	}
	downloadChanged.MarkManualChanges()

	uc.broadcastDownloadChanged(ctx, *downloadChanged)

	return nil
}
