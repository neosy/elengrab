package mediawatch

import "context"

func (uc *MediaWatch) ExecuteUpdateStats(ctx context.Context, workerID uint64) error {
	return uc.updateAllStats(ctx)
}
