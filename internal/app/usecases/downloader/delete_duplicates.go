package downloader

import (
	"context"
	"os"
	"path/filepath"
)

func (uc *YouTubeDownloader) DeleteDuplicates(ctx context.Context) error {
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

			err := uc.file.HardDelete(ctx, file.FileId)
			if err != nil {
				continue
			}

			delCnt++

			fPath := filepath.Join(uc.downloadsDir, file.FullName)
			if err := os.Remove(fPath); err != nil {
				uc.logger.Warn("Failed delete file", "filePath", fPath, "error", err)
				continue
			}
		}
		uc.logger.Info("Delete duplicate", "youtubeTitle", files[0].YoutubeTitle, "count", delCnt)
	}

	return nil
}
