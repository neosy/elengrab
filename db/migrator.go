package database

import (
	"embed"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type migratorInterface interface {
	newSource() (source.Driver, error)
	newMigrator(src source.Driver) (*migrate.Migrate, error)
}

type migratorFS struct {
	sqlDrv database.Driver
	fs     embed.FS
}

type migratorFile struct {
	sqlDrv database.Driver
	dir    string
}

func newMigratorFS(sqlDrv database.Driver, fs embed.FS) *migratorFS {
	return &migratorFS{
		sqlDrv: sqlDrv,
		fs:     fs,
	}
}

func (m *migratorFS) newSource() (source.Driver, error) {
	// Check if migrations directory exists
	drv, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migrations directory does not exist in embedded FS: %w", err)
	}

	return drv, nil
}

func (m *migratorFS) newMigrator(src source.Driver) (*migrate.Migrate, error) {
	// Create migrate instance pointing to the migrations folder
	migrator, err := migrate.NewWithInstance(
		"iofs", src,
		"sqlite", m.sqlDrv,
	)
	if err != nil {
		return nil, err
	}

	return migrator, nil
}

func newMigratorFile(sqlDrv database.Driver, dir string) *migratorFile {
	return &migratorFile{
		sqlDrv: sqlDrv,
		dir:    dir,
	}
}

func (m *migratorFile) newSource() (source.Driver, error) {
	if m.dir == "" {
		return nil, fmt.Errorf("migration directory is not specified")
	}

	// Check if migrations directory exists
	if _, err := os.Stat(m.dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory does not exist: %s", m.dir)
	}

	// Create file source driver
	drv, err := (&file.File{}).Open("file://" + m.dir)
	if err != nil {
		return nil, err
	}

	return drv, nil
}

func (m *migratorFile) newMigrator(src source.Driver) (*migrate.Migrate, error) {
	// Create migrate instance pointing to the migrations folder
	migrator, err := migrate.NewWithInstance(
		"file", src,
		"sqlite", m.sqlDrv,
	)
	if err != nil {
		return nil, err
	}

	return migrator, nil
}
