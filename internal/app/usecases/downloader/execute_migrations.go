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
		d.file,
		d.dlDataMigration,
	)

	return migrations.ExecuteMigrations(ctx)
}
