package taskexecutor

import (
	"context"

	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type Broadcaster struct {
	MediaAdded              func(ctx context.Context, userID uuid.UUID, download *ddownload.MediaDownload)
	MediaUpdated            func(ctx context.Context, downloadID uuid.UUID)
	DownloadProgressUpdated func(ctx context.Context, downloadID uuid.UUID)
	StartRefreshing         func(ctx context.Context, downloadID uuid.UUID)
	DownloadsCountUpdated   func(ctx context.Context, userID uuid.UUID)
}
