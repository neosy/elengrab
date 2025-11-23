package fileuc

import "context"

func (uc *File) ResetStatus(ctx context.Context) error {
	err := uc.FileRep.UpdateStatusToNew(ctx)
	if err != nil {
		uc.logger.Warn("Failed update status to new", "error", err)
		return err
	}
	return nil
}
