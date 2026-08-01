package downloader

import (
	"context"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ierrors "github.com/neosy/elengrab/internal/errors"
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

	return uc.mediaWatch.CreateMediaWatchEvent(&req, mediaDuration, uc)
}
