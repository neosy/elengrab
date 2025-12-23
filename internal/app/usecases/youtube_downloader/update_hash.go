package ytdownloader

import (
	"context"
	"os"
	"path"

	"github.com/neosy/elengrab/internal/app/utils"
)

func (uc *YouTubeDownloader) UpdateHash(ctx context.Context) error {
	files, err := uc.file.GetWithoutPartialHash(ctx)
	if err != nil {
		return err
	}

	for _, file := range files {
		fPath := path.Join(uc.downloadsDir, file.FullName)
		if _, err := os.Stat(fPath); err != nil {
			uc.logger.Warn("File not found", "filePath", fPath, "error", err)
			continue
		}

		h, err := utils.HashPartialMedia(fPath)
		if err != nil {
			uc.logger.Warn("Failed get hash partial", "filePath", fPath, "error", err)
			continue
		}

		file.PartialHash = &h
		err = uc.file.Update(ctx, file)
		if err != nil {
			continue
		}
	}

	return nil
}
