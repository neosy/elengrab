package mediadownload

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (uc *MediaDownload) ResetDownloadStateCache(ctx context.Context) error {
	downloads := make([]*ddownload.MediaDownload, 0)

	err := uc.downloadRepo().WithStatus(dtypes.MediaDownloadIncompleteStatuses()...).IterateAll(
		ctx,
		func(d *ddownload.MediaDownload) error {
			downloads = append(downloads, d)
			return nil
		},
	)
	if err != nil {
		return err
	}

	err = uc.dlStateCache.Clear(ctx)
	if err != nil {
		return err
	}

	for _, d := range downloads {
		err := uc.dlStateCache.SaveByDownload(ctx, d)
		if err != nil {
			return err
		}
	}

	return nil
}
