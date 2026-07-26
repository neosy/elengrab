package watchevent

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (uc *MediaWatchEvent) IterateGetAll(ctx context.Context, fn func(*ddownload.MediaWatchEvent) error) error {
	return uc.eventRep.IterateGetAll(ctx, fn)
}
