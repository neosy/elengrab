package migrations

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/dependencies"
	dlmigration "github.com/neosy/elengrab/internal/app/usecases/migrations/internal/download_data_migration"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type migrationProvider interface {
	RunMigrations(ctx context.Context) error
}

type migrations struct {
	logger *slog.Logger

	downloadsStorage pstorage.DownloadsStorage

	usecases dependencies.Usecases
	services dependencies.Services
}

func NewMigrations(
	logger *slog.Logger,
	downloadMigrationRepo persistence.DownloadDataMigrationRepositoryFactory,
	deps Dependencies,
) Migrations {
	services := deps.Services

	migrations := &migrations{
		logger: logger,

		downloadsStorage: deps.DownloadsStorage,

		usecases: dependencies.Usecases{
			Downloader:    deps.Usecases.Downloader,
			MediaDownload: deps.Usecases.Downloader.MediaDownload(),
			MediaWatch:    deps.Usecases.Downloader.MediaWatch(),

			DownloadMigration: dlmigration.NewDownloadMigration(logger, downloadMigrationRepo),
			Thumbnail:         deps.Usecases.Thumbnail,
		},

		services: services,
	}

	return migrations
}
