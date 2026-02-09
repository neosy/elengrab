package database

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type migratorInterface interface {
	newMigrator(src source.Driver) (*migrate.Migrate, error)
	newSource(wrapper sourceDriverWrapper) (source.Driver, sourceCleanup, error)
}

type migratorFS struct {
	sqlDrv database.Driver
	fs     fs.FS
}

type migratorFile struct {
	sqlDrv database.Driver
	dir    string
}

func newMigratorFS(sqlDrv database.Driver, fs fs.FS) *migratorFS {
	return &migratorFS{
		sqlDrv: sqlDrv,
		fs:     fs,
	}
}

func (m *migratorFS) newSource(wrapper sourceDriverWrapper) (source.Driver, sourceCleanup, error) {
	// Check if migrations directory exists
	drv, err := iofs.New(m.fs, ".")
	if err != nil {
		return nil, nil, fmt.Errorf("migrations directory does not exist in embedded FS: %w", err)
	}

	// Wrap driver if necessary
	if wrapper != nil {
		return wrapper(drv)
	}

	return drv, nil, nil
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

func (m *migratorFile) newSource(wrapper sourceDriverWrapper) (source.Driver, sourceCleanup, error) {
	if m.dir == "" {
		return nil, nil, fmt.Errorf("migration directory is not specified")
	}

	// Check if migrations directory exists
	if _, err := os.Stat(m.dir); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("migrations directory does not exist: %s", m.dir)
	}

	// Create file source driver
	drv, err := (&file.File{}).Open("file://" + m.dir)
	if err != nil {
		return nil, nil, err
	}

	// Wrap driver if necessary
	if wrapper != nil {
		return wrapper(drv)
	}

	return drv, nil, nil
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
