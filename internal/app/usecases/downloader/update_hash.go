package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/utils/hash"
)

func (uc *downloader) UpdateHash(ctx context.Context) error {
	downloads, err := uc.download.GetWithoutPartialHash(ctx)
	if err != nil {
		return err
	}

	for _, download := range downloads {
		filePath := uc.downloadsStorage.Path(download.FileFullName)

		exists, err := uc.downloadsStorage.Exists(download.FileFullName)
		if err != nil || !exists {
			uc.logger.Warn("File not found", "filePath", filePath, "error", err)
			continue
		}

		h, err := hash.FilePartialHash(filePath)
		if err != nil {
			uc.logger.Warn("Failed get hash partial", "filePath", filePath, "error", err)
			continue
		}

		download.PartialHash = &h
		err = uc.download.Update(ctx, nil, download)
		if err != nil {
			continue
		}
	}

	return nil
}
