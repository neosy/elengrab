package migrations

import (
	"context"
	"fmt"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *migrations) fillTitleForInstagram(ctx context.Context) (bool, error) {
	medias, err := m.getAllDownloads(ctx, false,
		func(download *ddownload.MediaDownload) bool {
			if hostdetect.Instagram(download.MediaURL) {
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

	m.logger.Debug("Found instagram medias", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.logger.Debug("Extract media info", "index", i+1, "total", len(medias))

		info, err := m.services.downloader.FetchInfo(ctx, media.MediaURL)
		if err != nil {
			continue
		}

		if info.Title == "" || media.MediaTitle == info.Title {
			continue
		}

		err = m.usecases.download.Patch(
			ctx, nil,
			media.DownloadID,
			func(download *ddownload.MediaDownload) error {
				download.MediaTitle = info.Title
				if info.Description != "" && (download.MediaDescription == nil || *download.MediaDescription != info.Description) {
					download.MediaDescription = &info.Description
				}
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
