package downloader

import (
	"context"
	"net/http"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (uc *Downloader) PatchMediaDownload(
	ctx context.Context,
	authCtx dauth.UserContext,
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

	err = uc.validateDownloadWriteAccess(authCtx, download)
	if err != nil {
		return err
	}

	var needUpdate bool

	if req.MediaTitle != nil && *req.MediaTitle != download.MediaTitle {
		download.MediaTitle = *req.MediaTitle
		needUpdate = true
	}

	var (
		description     string
		origDescription string
	)

	if req.MediaDescription != nil {
		description = *req.MediaDescription
	}

	if download.MediaDescription != nil {
		origDescription = *download.MediaDescription
	}

	if description != origDescription {
		download.MediaDescription = &description
		needUpdate = true
	}

	if req.Visibility != nil && *req.Visibility != download.Visibility {
		download.Visibility = *req.Visibility
		needUpdate = true
	}

	if !needUpdate {
		return errorx.NewHTTPMessage("No changes to update", http.StatusBadRequest)
	}

	if err := download.Validate(); err != nil {
		return err
	}

	return uc.download.Update(ctx, &authCtx.UserID, download)
}
