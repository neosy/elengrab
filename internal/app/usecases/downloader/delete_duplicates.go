package downloader

import (
	"context"
)

// DeleteDuplicates deletes duplicate downloads based on their hash values.
func (uc *Downloader) DeleteDuplicates(ctx context.Context) error {
	rows, err := uc.download.GetDuplicateHashes(ctx, uc.deleteDuplicatesScope)
	if err != nil {
		return err
	}

	for _, r := range rows {
		downloads, err := uc.download.GetByPartialHash(ctx, r)
		if err != nil {
			continue
		}
		if len(downloads) < 2 {
			continue
		}

		var delCnt uint16
		for i, download := range downloads {
			if i == 0 {
				continue
			}

			err := uc.download.HardDelete(ctx, download.DownloadID)
			if err != nil {
				continue
			}
			uc.broadcastDownloadDelete(ctx, download)

			delCnt++

			uc.deleteThumbnails(ctx, download)

			if err := uc.downloadsStorage.Delete(download.FileFullName); err != nil {
				uc.logger.Warn("Failed delete file", "filePath", uc.downloadsStorage.Path(download.FileFullName), "error", err)
				continue
			}
		}
		uc.logger.Info("Delete duplicate", "mediaTitle", downloads[0].MediaTitle, "count", delCnt)
	}

	return nil
}
