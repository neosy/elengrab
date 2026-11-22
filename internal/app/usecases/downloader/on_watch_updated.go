package downloader

import (
	"context"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *downloader) onWatchStatsUpdated(ctx context.Context, key dto.MediaWatchUpdateKey) {
	authCtx, err := uc.authContextFromUserAndSession(ctx, key.UserID, key.SessionID)
	if err != nil {
		return
	}

	uc.broadcastWatchStatsUpdated(ctx, authCtx, key.DownloadID)
}

func (uc *downloader) onWatchStatsUpdatedDependencies(ctx context.Context, stat *ddownload.MediaWatchStat) {
	if stat == nil {
		return
	}
	uc.searchIndex.UpdateViews(ctx, stat.DownloadID, stat.Views)
}

func (uc *downloader) onWatchUserPositionUpdated(ctx context.Context, key dto.MediaWatchUpdateKey) {
	authCtx, err := uc.authContextFromUserAndSession(ctx, key.UserID, key.SessionID)
	if err != nil {
		return
	}

	uc.broadcastWatchPositionUpdatedToAuth(ctx, authCtx, key.DownloadID)
}

func (uc *downloader) authContextFromUserAndSession(ctx context.Context, userID, sessionID uuid.UUID) (dauth.AuthContext, error) {
	authCtx := dauth.NewEmptyAuthContext()

	if userID == dauth.AnonymousUserID() {
		return dauth.NewAnonymousSessionAuthContext(sessionID), nil
	}

	user, err := uc.authSrv.GetByUserID(ctx, userID)
	if err != nil {
		return dauth.NewEmptyAuthContext(), err
	}

	err = authCtx.InitFromUser(user)
	if err != nil {
		return dauth.NewEmptyAuthContext(), err
	}

	return authCtx, nil
}
