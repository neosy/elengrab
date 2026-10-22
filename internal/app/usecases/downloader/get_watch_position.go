package downloader

import (
	"context"
	"time"

	"uuid"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *downloader) GetLastWatchPosition(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) (time.Duration, error) {
	return uc.mediaWatch.GetLastUserWatchPosition(ctx, downloadID, authCtx.UserID, authCtx.AnonSessionID)
}
