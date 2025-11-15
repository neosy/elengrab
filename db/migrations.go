package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
)

//go:embed migrations/*
var migrationsFS embed.FS

type MigrationConfig struct {
	// path to migration files
	Dir string
}

// Migrations represents a database migration manager.
type Migrations struct {
	db     *sql.DB
	config *MigrationConfig
}

// NewMigrations creates a new Migrations instance.
//
// Parameters:
//   - db: *sql.DB — the database connection to apply migrations on.
//
// Returns:
//   - *Migrations: a new Migrations manager instance.
func NewMigrations(db *sql.DB, config *MigrationConfig) *Migrations {
	return &Migrations{
		db:     db,
		config: config,
	}
}

func (m *Migrations) applyUp(migrator *migrate.Migrate) error {
	// Apply all up migrations
	err := migrator.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("All migrations are already applied. Database is up to date.")
			return nil
		}
		return fmt.Errorf("migration failed: %v", err)
	}
	log.Println("New migrations applied successfully.")
	return nil
}

// ApplyMigrations applies all available migrations.
//
// Returns:
//   - error: non-nil if any step fails.
func (m *Migrations) ApplyMigrations() error {
	// Create SQLite driver instance for golang-migrate
	sqlDrv, err := sqlite.WithInstance(m.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	var migrator migratorInterface
	if m.config != nil && m.config.Dir != "" {
		migrator = newMigratorFile(sqlDrv, m.config.Dir)
	}

	if migrator == nil {
		migrator = newMigratorFS(sqlDrv, migrationsFS)
	}

	src, err := migrator.newSource()
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
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
