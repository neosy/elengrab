package dlmigration

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DataMigration struct {
	logger *slog.Logger

	// repositories
	dataMigrationRep persistence.DownloadDataMigrationRepository
}

func NewdataMigration(
	logger *slog.Logger,
	dataMigrationRep persistence.DownloadDataMigrationRepository,
) *DataMigration {
	return &DataMigration{
		logger:           logger,
		dataMigrationRep: dataMigrationRep,
	}
}
