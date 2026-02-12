package database

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
)

// SQLiteDriverSource creates a new source driver for SQLite.
func newSQLiteDriverSource(sqliteDir string) sourceDriverWrapper {
	return func(drv source.Driver) (source.Driver, sourceCleanup, error) {
		if sqliteDir == "" {
			return nil, nil, errors.New("SQLite directory is not set")
		}

		// Create a temporary directory to store migrations files
		tmpDir, err := os.MkdirTemp("", "migrations_tmp_*")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create temp dir: %w", err)
		}

		cleanup := func() error {
			if tmpDir == "" {
				return nil
			}
			return os.RemoveAll(tmpDir)
		}

		// Iterate over all UP migrations
		version, err := drv.First()
		if err != nil {
			cleanup()
			if err == os.ErrNotExist {
				// No migrations, nothing to copy
				return drv, nil, nil
			}
			return nil, nil, fmt.Errorf("failed to get first migration version: %w", err)
		}

		for {
			rc, id, err := drv.ReadUp(version)
			if err != nil {
				cleanup()
				return nil, nil, fmt.Errorf("failed to read UP migration %d: %w", version, err)
			}

			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				cleanup()
				return nil, nil, fmt.Errorf("failed to read content of migration %d: %w", version, err)
			}

			// Replace oldPath with newPath
			newContent := strings.ReplaceAll(string(content), "${SOURCE_DB}", sqliteDir)

			// Write to tmpDir
			fileName := fmt.Sprintf("%03d_%s.up.sql", version, id)
			outFile := filepath.Join(tmpDir, fileName)
			if err := os.WriteFile(outFile, []byte(newContent), 0o644); err != nil {
				cleanup()
				return nil, nil, fmt.Errorf("failed to write migration %s to tmpDir: %w", id, err)
			}

			// ---- DOWN ----
			if rc, id, err := drv.ReadDown(version); err == nil {
				content, _ := io.ReadAll(rc)
				rc.Close()

				fileName := fmt.Sprintf("%03d_%s.down.sql", version, id)
				outFile := filepath.Join(tmpDir, fileName)
				if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
					cleanup()
					return nil, nil, fmt.Errorf("failed to write migration %s to tmpDir: %w", id, err)
				}

			}

			// Move to next version
			nextVersion, err := drv.Next(version)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					break // no more migrations
				}
				cleanup()
				return nil, nil, fmt.Errorf("failed to get next version after %d: %w", version, err)
			}
			version = nextVersion
		}

		// Create file source driver
		sourceURL := filePathForMigrate(tmpDir)
		drv, err = (&file.File{}).Open(sourceURL)
		if err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to open file source %s: %w", sourceURL, err)
		}

		return drv, cleanup, nil
	}
}
