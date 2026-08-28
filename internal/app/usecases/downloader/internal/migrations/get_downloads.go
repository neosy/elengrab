package migrations

import (
	"context"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *migrations) getAllDownloads(
	ctx context.Context, includeDeleted bool,
	match func(dowMediaDownloadnload *ddownload.MediaDownload) bool,
) ([]*ddownload.MediaDownload, error) {
	var downloads []*ddownload.MediaDownload

	iterateGetAll := m.usecases.download.IterateGetAll
	if includeDeleted {
		iterateGetAll = m.usecases.download.IterateGetAllWithDeleted
	}

	err := iterateGetAll(ctx,
		func(download *ddownload.MediaDownload) error {
			if download == nil || download.MediaInfo == nil {
				return nil
			}

			if download.MediaURL == "" {
				return nil
			}

			if !match(download) {
				return nil
			}

			downloads = append(downloads, download)

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return downloads, nil
}
