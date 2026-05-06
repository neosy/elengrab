package migrations

import (
	"context"
	"fmt"
	"path/filepath"

	ddownload "github.com/neosy/elengrab/internal/domain/download"
	nfile "github.com/neosy/elengrab/internal/pkg/file"
)

func (m *migrations) migrateMoveDownloadsToStorage(ctx context.Context) error {
	exists, err := m.dlMigration.Exists(ctx, migrateMoveDownloadsToStorageID)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	fileNames, err := m.dlFile.GetAllFullNames(ctx, true)
	if err != nil {
		return err
	}

	markMigration := func() error {
		migration := &ddownload.DataMigration{
			MigrationID: migrateMoveDownloadsToStorageID,
		}

		err = m.dlMigration.Insert(ctx, migration)
		if err != nil {
			return err
		}

		return nil
	}

	if len(fileNames) == 0 {
		return markMigration()
	}

	var hasErr = false
	for fName := range fileNames {
		filePath := filepath.Join(m.dlStorage.BasePath(), fName)
		exists, err := nfile.FileExists(filePath)
		if err != nil {
			m.logger.Warn("Failed to exists file", "filePath", filePath, "error", err)
			hasErr = true
		}
		if err == nil && !exists {
			continue
		}

		err = m.dlStorage.Move(filePath, fName)
		if err != nil {
			m.logger.Warn("Failed to move storage", "filePath", filePath, "error", err)
			hasErr = true
			continue
		}
	}

	if hasErr {
		return fmt.Errorf("errors in the migration process '%s'", migrateMoveDownloadsToStorageID)
	}

	return markMigration()
}
