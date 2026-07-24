package pworkers

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type WatchEventRunner interface {
	ExecuteCreateMediaWatchEvent(ctx context.Context, workerID uint64, req *dto.CreateMediaWatchEventRequest) error
	ExecuteUpdateStats(ctx context.Context, workerID uint64) error
}
