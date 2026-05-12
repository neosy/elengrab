package dlmigration

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadMigration struct {
	logger *slog.Logger

	// repositories
	dataMigrationRep persistence.DownloadDataMigrationRepository
}

func NewDownloadMigration(
	logger *slog.Logger,
	dataMigrationRep persistence.DownloadDataMigrationRepository,
) *DownloadMigration {
	return &DownloadMigration{
		logger:           logger,
		dataMigrationRep: dataMigrationRep,
	}
}
