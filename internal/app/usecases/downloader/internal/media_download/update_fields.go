package mediadownload

import (
	"context"
	"net/http"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (uc *MediaDownload) UpdateFields(
	ctx context.Context,
	req dto.PatchMediaDownloadRequest,
) error {
	var download *ddownload.MediaDownload

	update := func(ctx context.Context) error {
		var err error

		download, err = uc.GetByDownloadIDNoCache(ctx, req.DownloadID)
		if err != nil {
			return err
		}

		var needUpdate bool

		if req.MediaTitle != nil && !helper.ValuesEqual(&download.MediaTitle, req.MediaTitle) {
			download.MediaTitle = *req.MediaTitle
			needUpdate = true
		}

		var reqMediaDescription *string
		if req.MediaDescription != nil {
			reqMediaDescription = *req.MediaDescription
		}

		if !helper.ValuesEqual(download.MediaDescription, reqMediaDescription) {
			download.MediaDescription = reqMediaDescription
			needUpdate = true
		}

		if req.Visibility != nil && !helper.ValuesEqual(&download.Visibility, req.Visibility) {
			download.Visibility = *req.Visibility
			needUpdate = true
		}

		if !needUpdate {
			return errorx.NewHTTPMessage("No changes to update", http.StatusBadRequest)
		}

		if err := download.Validate(); err != nil {
			return err
		}

		download.NormalizeForSave()

		err = uc.downloadRepo().Update(ctx, download)
		if err != nil {
			uc.logger.Warn("Update record error", "error", err)
			return err
		}

		return nil
	}

	err := uc.Tx(ctx, update)
	if err != nil {
		return err
	}

	uc.dlStateCache.PatchDownload(ctx, req)
	uc.downloadCacheRep.Delete(ctx, req.DownloadID)

	uc.updateDependencies(ctx, download)

	return nil
}
