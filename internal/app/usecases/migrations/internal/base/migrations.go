package base

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/dependencies"
	"github.com/neosy/elengrab/internal/app/usecases/migrations/internal/registry"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Migrations struct {
	logger *slog.Logger

	migrationList registry.MigrationList

	dlStorage pstorage.DownloadsStorage

	usecases dependencies.Usecases
	services dependencies.Services
}

func NewMigrations(
	logger *slog.Logger,
	dlStorage pstorage.DownloadsStorage,
	usecases dependencies.Usecases,
	services dependencies.Services,
) Migrations {
	migrations := Migrations{
		logger: logger,

		migrationList: registry.NewMigrationList(),

		dlStorage: dlStorage,

		usecases: usecases,
		services: services,
	}

	return migrations
}

func (m *Migrations) MarkMigration(ctx context.Context, migrationID string) error {
	return m.usecases.DownloadMigration.MarkMigration(ctx, migrationID)
}

func (m *Migrations) Add(id string, run registry.MigrationRunner) {
	m.migrationList.Add(id, run)
}

func (m *Migrations) Logger() *slog.Logger {
	return m.logger
}

func (m *Migrations) Usecases() dependencies.Usecases {
	return m.usecases
}

func (m *Migrations) Services() dependencies.Services {
	return m.services
}

func (m *Migrations) DownloadsStorage() pstorage.DownloadsStorage {
	return m.dlStorage
}
