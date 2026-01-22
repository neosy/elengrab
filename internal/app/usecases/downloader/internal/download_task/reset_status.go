package dltask

import "context"

func (uc *DownloadTask) ResetStatus(ctx context.Context) error {
	err := uc.TaskRep.UpdateStatusToNew(ctx)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}
	return nil
}
