package dlmigration

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadMigration struct {
	logger *slog.Logger

	// repositories
	dataMigrationRepo persistence.DownloadDataMigrationRepositoryFactory
}

func NewDownloadMigration(
	logger *slog.Logger,
	dataMigrationRepo persistence.DownloadDataMigrationRepositoryFactory,
) *DownloadMigration {
	return &DownloadMigration{
		logger:           logger,
		dataMigrationRepo: dataMigrationRepo,
	}
}
