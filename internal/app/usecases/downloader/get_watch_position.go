package downloader

import (
	"context"
	"time"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *downloader) GetLastWatchPosition(
	ctx context.Context,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
) (time.Duration, error) {
	return uc.mediaWatch.GetLastUserWatchPosition(ctx, downloadID, authCtx.UserID, authCtx.AnonSessionID)
}
