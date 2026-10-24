package pworkers

import (
	"context"
	"time"
)

// DownloadMaintenanceRunner defines the interface for a download maintenance runner.
type DownloadMaintenanceRunner interface {
	UpdateHash(ctx context.Context) error
	DeleteDuplicates(ctx context.Context) error
	DeleteMissingMediaDownloads(ctx context.Context, enableMoveUnmatchedFiles bool) error
	DeleteFailedMediaDownloads(ctx context.Context) error
	ResolveInconsistencies(ctx context.Context, orphanGracePeriod time.Duration) error
}
