package uploadexecutor

import (
	"context"
	"uuid"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type Broadcaster struct {
	DownloadAdded         func(ctx context.Context, userID uuid.UUID, download *ddownload.MediaDownload)
	DownloadUpdated       func(ctx context.Context, downloadID uuid.UUID)
	StartRefreshing       func(ctx context.Context, downloadID uuid.UUID)
	DownloadsCountUpdated func(ctx context.Context, userID uuid.UUID)
}
