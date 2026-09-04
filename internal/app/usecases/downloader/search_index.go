package downloader

import (
	"context"

	"github.com/google/uuid"
	mediadownload "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_download"
	mediawatch "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch"
	searchindex "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/search_index"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
)

type SearchIndex interface {
	Build(ctx context.Context) error
}

type searchIndex struct {
	*searchindex.SearchIndex
	download   *mediadownload.MediaDownload
	mediaWatch *mediawatch.MediaWatch
}

func (u *downloader) SearchIndex() SearchIndex {
	return &searchIndex{
		SearchIndex: u.searchIndex,
		download:    u.download,
		mediaWatch:  u.mediaWatch,
	}
}
func (u *searchIndex) CreateMediaDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	return u.SearchIndex.CreateMediaDownload(ctx, download)
}

func (u *searchIndex) SaveMediaDownload(ctx context.Context, download *ddownload.MediaDownload) error {
	return u.SearchIndex.SaveMediaDownload(ctx, download)
}

func (u *searchIndex) UpdateViews(ctx context.Context, downloadID uuid.UUID, views uint32) error {
	return u.SearchIndex.UpdateViews(ctx, downloadID, views)
}

func (u *searchIndex) Build(ctx context.Context) error {
	return u.download.IterateGetAll(
		ctx,
		func(download *ddownload.MediaDownload) error {
			if download == nil {
				return nil
			}

			err := u.SearchIndex.CreateMediaDownload(ctx, download)
			if err != nil {
				return err
			}

			views, err := u.mediaWatch.GetViews(ctx, download.DownloadID)
			if err != nil {
				return err
			}

			if views == 0 {
				return nil
			}

			err = u.SearchIndex.UpdateViews(ctx, download.DownloadID, views)
			if err != nil {
				return err
			}

			return nil
		},
	)
}
