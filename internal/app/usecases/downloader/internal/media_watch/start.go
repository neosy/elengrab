package mediawatch

import "context"

func (uc *MediaWatch) Start(ctx context.Context) {
	uc.startStatsUpdater(ctx)
}
