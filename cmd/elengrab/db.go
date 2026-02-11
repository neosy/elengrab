package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	database "github.com/neosy/elengrab/db"
	iconfig "github.com/neosy/elengrab/infrastructure/config"
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
	authDB *sql.DB,
	mainDB *sql.DB,
	mediaDB *sql.DB,
) error {
	sqliteDir := absPath(cfg.Elengrab.AppDir, cfg.SQLite.DataDir)

	migration := database.NewMigrations(logger, authDB, persistence.DBAuthName, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", persistence.DBAuthName), "error", err)
		return err
	}

	migration = database.NewMigrations(logger, mainDB, persistence.DBMainName, nil)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", persistence.DBMainName), "error", err)
		return err
	}

	migration = database.NewMigrations(
		logger,
		mediaDB,
		persistence.DBMediaName,
		&database.MigrationConfig{
			SQLiteDir:       sqliteDir,
			NoTxWrap:        true,
			ResolveSourceDB: true,
		},
	)
	if err := migration.ApplyMigrations(); err != nil {
		logger.Error(fmt.Sprintf("Failed to migration SQLite database: %v", persistence.DBMediaName), "error", err)
		return err
	}

	return nil
}
