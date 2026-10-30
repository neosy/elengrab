package uploadexecutor

import (
	"context"
	"uuid"
)

type Broadcaster struct {
	MediaUpdated    func(ctx context.Context, downloadID uuid.UUID)
	StartRefreshing func(ctx context.Context, downloadID uuid.UUID)
}
