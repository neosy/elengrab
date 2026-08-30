package maintenance

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

const (
	backupFileExt    = "sqlite"
	backupDateFormat = "2006-01-02_15-04-05"
)

type maintenance struct {
	logger *slog.Logger

	// Repositories
	repositories persistence.Repositories

	// options
	appName             string
	databaseBackupsDir  string
	databaseBackupsKeep int
}

func NewMaintenance(
	logger *slog.Logger,

	// Repositories
	repositories persistence.Repositories,

	// options
	appName string,
	databaseBackupsDir string,
	databaseBackupsKeep int,
) *maintenance {
	return &maintenance{
		logger:       logger,
		repositories: repositories,

		// options
		appName:             appName,
		databaseBackupsDir:  databaseBackupsDir,
		databaseBackupsKeep: databaseBackupsKeep,
	}
}
