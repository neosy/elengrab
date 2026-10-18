package processingtask

import "context"

func (uc *ProcessingTask) ResetStatus(ctx context.Context) error {
	err := uc.TaskRepo().UpdateStatusToNew(ctx)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}
	return nil
}
