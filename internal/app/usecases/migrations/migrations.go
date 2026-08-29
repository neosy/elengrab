package migrations

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/deferred"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/dependencies"
	dlmigration "github.com/neosy/elengrab/internal/app/usecases/migrations/internal/download_data_migration"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/required"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Usecases struct {
	Downloader downloader.InternalDownloader
	Thumbnail  *thumbnail.Thumbnail
}
type Services = dependencies.Services

type Dependencies struct {
	DownloadsStorage pstorage.DownloadsStorage

	Usecases Usecases
	Services Services
}

type migrationProvider interface {
	RunMigrations(ctx context.Context) error
}

type Migrations interface {
	// RunRequiredMigrations executes all required migrations.
	// The application cannot continue until these migrations complete successfully.
	RunRequiredMigrations(ctx context.Context) error

	// RunDeferredMigrations executes deferred migrations.
	// These migrations are not required for startup and can be run after the application is ready.
	RunDeferredMigrations(ctx context.Context) error
}

type migrations struct {
	logger *slog.Logger

	required migrationProvider
	deferred migrationProvider

	// Usecases
	downloadMigration *dlmigration.DownloadMigration
}

func NewMigrations(
	logger *slog.Logger,
	downloadMigrationRepo persistence.DownloadDataMigrationRepositoryFactory,
	deps Dependencies,
) Migrations {
	usecases := dependencies.Usecases{
		Downloader:    deps.Usecases.Downloader,
		MediaDownload: deps.Usecases.Downloader.MediaDownload(),
		MediaWatch:    deps.Usecases.Downloader.MediaWatch(),

		DownloadMigration: dlmigration.NewDownloadMigration(logger, downloadMigrationRepo),
		Thumbnail:         deps.Usecases.Thumbnail,
	}

	services := deps.Services

	migrations := &migrations{
		logger: logger,

		required: required.NewMigrations(logger, deps.DownloadsStorage, usecases, services),
		deferred: deferred.NewMigrations(logger, deps.DownloadsStorage, usecases, services),
	}

	return migrations
}
