package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	database "github.com/neosy/elengrab/db"
	iconfig "github.com/neosy/elengrab/internal/config"
	"github.com/neosy/elengrab/internal/ports/persistence"
	sqliterep "github.com/neosy/elengrab/internal/repository/sqlite"
)

// newDB initializes a new SQLite database connection.
func newDB(logger *slog.Logger, dbPath string) (*sql.DB, error) {
	db, err := sqliterep.InitDB(logger, dbPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize SQLite database: %v", filepath.Base(dbPath)), "error", err)
		return nil, err
	}

	return db, nil
}

// applyMigrations applies the necessary database migrations for the SQLite databases.
func applyMigrations(
	logger *slog.Logger,
	cfg *iconfig.Config,
	authEntry persistence.DBEntry,
	mainEntry persistence.DBEntry,
	mediaEntry persistence.DBEntry,
	linkEntry persistence.DBEntry,
	watchEventEntry persistence.DBEntry,
	searchIndexEntry persistence.DBEntry,
) error {
	sqliteDir := absPath(cfg.Elengrab.RootDir, cfg.SQLite.DataDir)

	migration := database.NewMigrations(logger, authEntry, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", authEntry.DBName()), "error", err)
		return err
	}

	migration = database.NewMigrations(logger, mainEntry, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", mainEntry.DBName()), "error", err)
		return err
	}

	migration = database.NewMigrations(
		logger,
		mediaEntry,
		&database.MigrationConfig{
			SQLiteDir:       sqliteDir,
			NoTxWrap:        true,
			ResolveSourceDB: true,
		},
	)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", mediaEntry.DBName()), "error", err)
		return err
	}

	migration = database.NewMigrations(logger, linkEntry, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", linkEntry.DBName()), "error", err)
		return err
	}

	migration = database.NewMigrations(logger, watchEventEntry, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", watchEventEntry.DBName()), "error", err)
		return err
	}

	migration = database.NewMigrations(logger, searchIndexEntry, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", searchIndexEntry.DBName()), "error", err)
		return err
	}

	return nil
}
