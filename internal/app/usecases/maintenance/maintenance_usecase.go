package maintenance

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

const (
	backupFileExt    = "sqlite"
	backupDateFormat = "2006-01-02_15-04-05"
)

type Maintenance struct {
	logger *slog.Logger

	// Database
	database persistence.Database

	// options
	appName             string
	databaseBackupsDir  string
	databaseBackupsKeep int
}

func NewMaintenance(
	logger *slog.Logger,

	// repositoies
	database persistence.Database,

	// options
	appName string,
	databaseBackupsDir string,
	databaseBackupsKeep int,
) *Maintenance {
	return &Maintenance{
		logger:   logger,
		database: database,

		// options
		appName:             appName,
		databaseBackupsDir:  databaseBackupsDir,
		databaseBackupsKeep: databaseBackupsKeep,
	}
}
