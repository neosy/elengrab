package migrations

import (
	"context"
	"fmt"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

func (m *migrations) fillMediaDuration(ctx context.Context) (bool, error) {
	medias, err := m.getAllDownloads(ctx, false,
		func(download *ddownload.MediaDownload) bool {
			if download.MediaInfo != nil && download.MediaInfo.DurationMs == 0 {
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

	m.logger.Debug("Found media with zero duration", "count", len(medias))

	var hasError bool

	for i, media := range medias {
		m.logger.Debug(
			"Extract media duration",
			"index", i+1,
			"total", len(medias),
			"downloadID", media.DownloadID.String(),
			"mediaTitle", media.MediaTitle,
			"format", media.MediaInfo.Format.String(),
		)

		filePath := m.dlStorage.Path(media.FileFullName)
		duration, err := m.services.ffmpeg.ExtractDurationMs(ctx, filePath)
		if err != nil {
			continue
		}

		if duration == 0 {
			continue
		}

		err = m.usecases.download.Patch(
			ctx, nil,
			media.DownloadID,
			func(download *ddownload.MediaDownload) bool {
				download.MediaInfo.SetDuration(duration)
				return true
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
