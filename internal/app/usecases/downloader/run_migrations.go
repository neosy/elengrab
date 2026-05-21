package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/migrations"
)

func (d *Downloader) RunRequiredMigrations(ctx context.Context) error {
	migrations := migrations.NewMigrations(
		d.logger,
		d.downloadsStorage,
		// usecases
		d.download,
		d.downloadMigration,
		d.thumbnail,
		// services
		d.downloaderSrv,
		d.ffmpegSrv,
	)
	return migrations.RunRequiredMigrations(ctx)
}

func (d *Downloader) RunDeferredMigrations(ctx context.Context) error {
	migrations := migrations.NewMigrations(
		d.logger,
		d.downloadsStorage,
		// usecases
		d.download,
		d.downloadMigration,
		d.thumbnail,
		// services
		d.downloaderSrv,
		d.ffmpegSrv,
	)
	return migrations.RunDeferredMigrations(ctx)
}
