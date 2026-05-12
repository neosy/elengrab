package downloader

import (
	"context"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/migrations"
)

const (
	migrateMoveDownloadsToStorageID = "migrate_move_downloads_to_storage"
)

func (d *Downloader) ExecuteMigrations(ctx context.Context) error {
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
	return migrations.ExecuteMigrations(ctx)
}
