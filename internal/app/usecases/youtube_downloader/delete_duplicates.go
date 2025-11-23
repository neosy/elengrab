package ytdownloader

import (
	"context"
	"os"
	"path"
)

func (uc *YouTubeDownloader) DeleteDuplicates(ctx context.Context) error {
	hashes, err := uc.file.GetDuplicateHashes(ctx)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		files, err := uc.file.GetByPartialHash(ctx, h)
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

			err := uc.file.Delete(ctx, file.FileId)
			if err != nil {
				continue
			}

			delCnt++

			fPath := path.Join(uc.downloadsDir, file.FullName)
			if err := os.Remove(fPath); err != nil {
				uc.logger.Warn("Failed delete file", "filePath", fPath, "error", err)
				continue
			}
		}
		uc.logger.Debug("Delete duplicate", "youtubeTitle", files[0].YoutubeTitle, "count", delCnt)
	}

	return nil
}
