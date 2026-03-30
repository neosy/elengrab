package downloader

import (
	"context"
	"os"
	"path/filepath"

	"github.com/neosy/elengrab/internal/app/utils"
)

func (uc *Downloader) UpdateHash(ctx context.Context) error {
	files, err := uc.file.GetWithoutPartialHash(ctx)
	if err != nil {
		return err
	}

	for _, file := range files {
		fPath := filepath.Join(uc.downloadsDir, file.FullName)
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
