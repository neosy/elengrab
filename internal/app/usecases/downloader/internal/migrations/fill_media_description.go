package migrations

import (
	"context"
	"fmt"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *migrations) fillMediaDescription(ctx context.Context) (bool, error) {
	medias, err := m.getAllDownloads(ctx, false,
		func(download *ddownload.MediaDownload) bool {
			if download.MediaDescription == nil || *download.MediaDescription == "" {
				return true
			}
			return false
		},
	)
	if err != nil {
		return false, err
	}

	if len(medias) == 0 {
		return true, nil
	}

	m.logger.Debug("Found medias with empty description", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.logger.Debug("Extract media info", "index", i+1, "total", len(medias))

		info, err := m.services.downloader.FetchInfo(ctx, media.MediaURL)
		if err != nil {
			continue
		}

		if info.Description == "" || (media.MediaDescription != nil && *media.MediaDescription == info.Description) {
			continue
		}

		err = m.usecases.download.Patch(
			ctx, nil,
			media.DownloadID,
			func(download *ddownload.MediaDownload) error {
				download.MediaDescription = &info.Description
				return nil
			},
		)
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return false, fmt.Errorf("errors in the migration process")
	}

	return true, nil
}
