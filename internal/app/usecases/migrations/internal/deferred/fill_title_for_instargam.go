package deferred

import (
	"context"
	"fmt"

	hostdetect "github.com/neosy/elengrab/internal/app/utils/host_detect"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *migrations) fillTitleForInstagram(ctx context.Context) error {
	medias, err := m.getAllDownloads(ctx, false,
		func(media *ddownload.MediaDownload) bool {
			if media.SourceURL() == "" {
				return false
			}

			if hostdetect.Instagram(media.SourceURL()) {
				return true
			}

			return false
		},
	)
	if err != nil {
		return err
	}

	if len(medias) == 0 {
		return nil
	}

	m.Logger().Debug("Found instagram medias", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		if media.SourceType != dtypes.MediaDownloadSourceTypeURL {
			continue
		}

		m.Logger().Debug("Extract media info", "index", i+1, "total", len(medias))

		info, err := m.Services().Downloader.FetchInfo(ctx, media.SourceURL())
		if err != nil {
			continue
		}

		if info.Title == "" || media.MediaTitle == info.Title {
			continue
		}

		err = m.Usecases().MediaDownload.Patch(
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
		return fmt.Errorf("errors in the migration process")
	}

	return nil
}
