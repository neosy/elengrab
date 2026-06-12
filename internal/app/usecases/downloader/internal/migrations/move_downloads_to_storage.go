package migrations

import (
	"context"
	"fmt"
	"path/filepath"

	nfile "github.com/neosy/elengrab/internal/pkg/filex"
)

func (m *migrations) moveDownloadsToStorage(ctx context.Context) (bool, error) {
	fileNames, err := m.usecases.download.GetAllFullNames(ctx, true)
	if err != nil {
		return false, err
	}

	if len(fileNames) == 0 {
		return false, nil
	}

	var hasErr = false
	for fName := range fileNames {
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

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
		return false, fmt.Errorf("errors in the migration process")
	}

	return true, nil
}
