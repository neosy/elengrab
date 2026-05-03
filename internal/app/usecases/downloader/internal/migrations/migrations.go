package migrations

import (
	"log/slog"

	dlmigration "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/dowload_data_migration"
	fileuc "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/file"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type migrations struct {
	logger *slog.Logger

	dlStorage pstorage.DownloadsStorage

	dlFile      *fileuc.File
	dlMigration *dlmigration.DataMigration
}

func NewMigrations(
	logger *slog.Logger,
	dlStorage pstorage.DownloadsStorage,
	dlFile *fileuc.File,
	dlMigration *dlmigration.DataMigration,
) *migrations {
	return &migrations{
		logger: logger,

		dlStorage: dlStorage,

		dlFile:      dlFile,
		dlMigration: dlMigration,
	}
}
