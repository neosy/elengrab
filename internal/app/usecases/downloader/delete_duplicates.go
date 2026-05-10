package downloader

import (
	"context"
)

// DeleteDuplicates deletes duplicate files based on their hash values.
func (uc *Downloader) DeleteDuplicates(ctx context.Context) error {
	rows, err := uc.file.GetDuplicateHashes(ctx, uc.deleteDuplicatesScope)
	if err != nil {
		return err
	}

	for _, r := range rows {
		files, err := uc.file.GetByPartialHash(ctx, r)
		if err != nil {
			continue
		}
		if len(files) < 2 {
			continue
		}

		var delCnt uint16
		for i, file := range files {
			if i == 0 {
				continue
			}

			err := uc.file.HardDelete(ctx, file.FileID)
			if err != nil {
				continue
			}
			uc.broadcastFileDelete(file.UserID, file.FileID)

			delCnt++

			uc.deleteThumbnails(ctx, file)

			if err := uc.downloadsStorage.Delete(file.FullFileName); err != nil {
				uc.logger.Warn("Failed delete file", "filePath", uc.downloadsStorage.Path(file.FullFileName), "error", err)
				continue
			}
		}
		uc.logger.Info("Delete duplicate", "mediaTitle", files[0].MediaTitle, "count", delCnt)
	}

	return nil
}
