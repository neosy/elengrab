package downloader

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *Downloader) UpdateVisibility(
	ctx context.Context,
	authCtx dauth.UserContext,
	downloadID uuid.UUID,
	visibility dtypes.MediaVisibility,
) error {
	err := uc.validateWriteOperation(authCtx)
	if err != nil {
		return err
	}

	download, err := uc.download.GetByDownloadID(ctx, downloadID)
	if err != nil {
		return err
	}

	err = uc.validateDownloadWriteAccess(authCtx, download)
	if err != nil {
		return err
	}

	patch := func(d *ddownload.MediaDownload) bool {
		d.Visibility = visibility
		return true
	}

	return uc.download.Patch(ctx, &authCtx.UserID, downloadID, patch)
}
