package downloadstatus

import (
	"context"

	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type StatusUpdater interface {
	New(ctx context.Context, downloadID uuid.UUID) error

	Pending(
		ctx context.Context,
		downloadID uuid.UUID,
		taskId uuid.UUID,
		jobID uuid.UUID,
	) error

	Failed(
		ctx context.Context,
		downloadID uuid.UUID,
		mutate func(*ddownload.MediaDownload) error,
		message string,
	) error

	Working(
		ctx context.Context,
		downloadID uuid.UUID,
		taskId uuid.UUID,
		workerId uint64,
	) error

	Done(
		ctx context.Context,
		downloadID uuid.UUID,
		mutate func(*ddownload.MediaDownload) error,
	) error

	Refreshing(ctx context.Context, downloadID uuid.UUID) error
}
