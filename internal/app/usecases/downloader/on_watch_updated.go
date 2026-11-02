package downloader

import (
	"context"

	"uuid"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *downloader) onWatchUserStatsUpdated(ctx context.Context, _ dauth.AuthContext, downloadID uuid.UUID) {
	uc.broadcastWatchStatsUpdated(ctx, downloadID)
}

func (uc *downloader) onWatchStatsUpdated(ctx context.Context, stat *ddownload.MediaWatchStat) {
	if stat == nil {
		return
	}
	uc.searchIndex.UpdateViews(ctx, stat.DownloadID, stat.Views)
}

func (uc *downloader) onWatchUserPositionUpdated(ctx context.Context, authCtx dauth.AuthContext, downloadID uuid.UUID) {
	uc.broadcastWatchPositionUpdatedToAuth(ctx, authCtx, downloadID)
}
