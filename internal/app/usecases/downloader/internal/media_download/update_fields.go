package mediadownload

import (
	"context"
	"net/http"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (uc *MediaDownload) UpdateFields(
	ctx context.Context,
	req dto.PatchMediaDownloadRequest,
) error {
	return uc.Tx(ctx, func(ctx context.Context) error {
		var err error

		download, err := uc.GetByDownloadIDNoCache(ctx, req.DownloadID)
		if err != nil {
			return err
		}

		var needUpdate bool

		if req.MediaTitle != nil && !helper.ValuesEqual(&download.MediaTitle, req.MediaTitle) {
			download.MediaTitle = *req.MediaTitle
			needUpdate = true
		}

		if !helper.ValuesEqual(download.MediaDescription, req.MediaDescription) {
			download.MediaDescription = req.MediaDescription
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

		err = uc.downloadRep.Update(ctx, download)
		if err != nil {
			uc.logger.Warn("Update record error", "error", err)
			return err
		}

		uc.downloadCacheRep.Delete(ctx, download.DownloadID)

		uc.dlStateCache.Transaction(func(ctx context.Context) error {
			state, _ := uc.dlStateCache.FindByDownloadID(ctx, req.DownloadID)
			if state != nil && state.Download != nil {
				if req.MediaTitle != nil {
					state.Download.MediaTitle = *req.MediaTitle
				}
				if req.MediaDescription != nil {
					state.Download.MediaDescription = req.MediaDescription
				}
				if req.Visibility != nil {
					state.Download.Visibility = *req.Visibility
				}
				uc.dlStateCache.Save(ctx, state)
			}

			return nil
		})

		return nil
	})
}
