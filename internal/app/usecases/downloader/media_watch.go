package downloader

import (
	"context"

	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
)

type MediaWatch interface {
	RebuildUserChunks(ctx context.Context) error
	RebuildWatchStats(ctx context.Context) error
}

type mediaWatch struct {
	*mediawatch.MediaWatch
	download *mediadownload.MediaDownload
}

func (m *mediaWatch) RebuildUserChunks(ctx context.Context) error {
	return m.MediaWatch.RebuildUserChunks(ctx, m.download.FindByDownloadID)
}

func (m *mediaWatch) RebuildWatchStats(ctx context.Context) error {
	return m.MediaWatch.RebuildWatchStats(ctx, m.download.FindByDownloadID)
}

func (uc *downloader) MediaWatch() MediaWatch {
	return &mediaWatch{
		MediaWatch: uc.mediaWatch,
		download:   uc.download,
	}
}
