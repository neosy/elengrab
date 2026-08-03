package downloader

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (uc *Downloader) TrackMediaWatchEvent(
	ctx context.Context,
	authCtx dauth.AuthContext,
	req dto.TrackMediaWatchEventRequest,
) error {
	if req.DownloadID == uuid.Nil {
		uc.logger.Warn("Id for the DownloadID field is not defined")
		return apperrors.ErrDownloadIDIsNil
	}

	if err := req.Validate(); err != nil {
		return errorx.Errorf(
			"invalid track media watch event request: %w",
			err, errorx.WithHttpStatus(http.StatusBadRequest),
		)
	}

	mediaDownload, err := uc.download.FindByDownloadID(ctx, req.DownloadID)
	if err != nil {
		uc.logger.Error("Failed get media info", "error", err)
		return err
	}
	if mediaDownload == nil {
		return apperrors.ErrDownloadNotFound
	}

	if !uc.authz.HasMediaViewAccess(authCtx, mediaDownload) {
		return ierrors.ErrAccessDenied
	}

	var mediaDuration time.Duration
	if mediaDownload.MediaInfo != nil {
		mediaDuration = time.Duration(mediaDownload.MediaInfo.DurationMs) * time.Millisecond
	}

	if mediaDuration <= 0 {
		uc.logger.Warn("Media duration is not defined or is zero", "downloadID", req.DownloadID)
		return errorx.NewHTTP("Media duration is not available", http.StatusConflict)
	}

	return uc.mediaWatch.CreateMediaWatchEvent(&req, mediaDuration, uc)
}
