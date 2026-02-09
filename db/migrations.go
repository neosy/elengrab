package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

//go:embed download/migrations/*
var migrationsMainFS embed.FS

//go:embed auth/migrations/*
var migrationsAuthFS embed.FS

//go:embed media/migrations/*
var migrationsMediaFS embed.FS

func migrationsMainRoot() (fs.FS, error) {
	return fs.Sub(migrationsMainFS, "download/migrations")
}

func migrationsAuthRoot() (fs.FS, error) {
	return fs.Sub(migrationsAuthFS, "auth/migrations")
}

func migrationsMediaRoot() (fs.FS, error) {
	return fs.Sub(migrationsMediaFS, "media/migrations")
}

type MigrationConfig struct {
	// Path to SQLite database directory
	SQLiteDir string
	// path to migration files
	Dir string
	// Disable automatic transaction wrapping for migrations (required for ATTACH, etc.)
	NoTxWrap bool
	// Enable copying migrations and resolving ${SOURCE_DB} placeholders
	ResolveSourceDB bool
}

// Migrations represents a database migration manager.
type Migrations struct {
	logger *slog.Logger
	db     *sql.DB
	name   persistence.DBName
	fs     fs.FS
	config *MigrationConfig
}

// NewMigrations creates a new Migrations instance.
//
// Parameters:
//   - db: *sql.DB — the database connection to apply migrations on.
//
// Returns:
//   - *Migrations: a new Migrations manager instance.
func NewMigrations(logger *slog.Logger, db *sql.DB, dbName persistence.DBName, config *MigrationConfig) *Migrations {
	var fs fs.FS

	switch dbName {
	case persistence.DBMainName:
		fs, _ = migrationsMainRoot()
	case persistence.DBAuthName:
		fs, _ = migrationsAuthRoot()
	case persistence.DBMediaName:
		fs, _ = migrationsMediaRoot()
	}

	return &Migrations{
		logger: logger,
		db:     db,
		name:   dbName,
		fs:     fs,
		config: config,
	}
}

// applyUp applies all up migrations using the provided migrator.
func (m *Migrations) applyUp(migrator *migrate.Migrate) error {
	// Apply all up migrations
	err := migrator.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Debug(fmt.Sprintf("All migrations are already applied. Database '%s' is up to date.", m.name))
			return nil
		}
		return fmt.Errorf("migration db '%s' failed: %v", m.name, err)
	}
	m.logger.Info(fmt.Sprintf("New migrations db '%s' applied successfully.", m.name))
	return nil
}

// ApplyMigrations applies all available migrations.
//
// Returns:
//   - error: non-nil if any step fails.
func (m *Migrations) ApplyMigrations() error {
	sqliteConfig := &sqlite.Config{}
	if m.config != nil {
		sqliteConfig.NoTxWrap = m.config.NoTxWrap
	}

	// Create SQLite driver instance for golang-migrate
	sqlDrv, err := sqlite.WithInstance(m.db, sqliteConfig)
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	var migrator migratorInterface
	if m.config != nil && m.config.Dir != "" {
		migrator = newMigratorFile(sqlDrv, m.config.Dir)
	}

	if migrator == nil {
		migrator = newMigratorFS(sqlDrv, m.fs)
	}

	var driver sourceDriverWrapper = nil
	if m.config != nil && m.config.ResolveSourceDB {
		driver = newSQLiteDriverSource(m.config.SQLiteDir)
	}

	src, cleanup, err := migrator.newSource(driver)
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	mgr, err := migrator.newMigrator(src)
	if err != nil {
		return fmt.Errorf("failed to create migrator instance: %w", err)
	}

	// ******** rollback migration
	// if err := migrator.Force(2); err != nil {
	// 	log.Fatalf("failed to force version: %v", err)
	// }
	// if err := migrator.Steps(-1); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatalf("failed to rollback migration: %v", err)
	// }

	return m.applyUp(mgr)
}
