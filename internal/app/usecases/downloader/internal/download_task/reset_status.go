package dltask

import "context"

func (uc *DownloadTask) ResetStatus(ctx context.Context) error {
	err := uc.TaskRepo().UpdateStatusToNew(ctx)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}
	return nil
}
