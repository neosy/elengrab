package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/utils/hash"
)

func (uc *Downloader) UpdateHash(ctx context.Context) error {
	files, err := uc.file.GetWithoutPartialHash(ctx)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := uc.downloadsStorage.Path(file.FileFullName)

		exists, err := uc.downloadsStorage.Exists(file.FileFullName)
		if err != nil || !exists {
			uc.logger.Warn("File not found", "filePath", filePath, "error", err)
			continue
		}

		h, err := hash.FilePartialHash(filePath)
		if err != nil {
			uc.logger.Warn("Failed get hash partial", "filePath", filePath, "error", err)
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
