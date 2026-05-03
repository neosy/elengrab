package migrations

import (
	"context"
	"fmt"
)

func (m *migrations) ExecuteMigrations(ctx context.Context) error {
	var hasError = false

	err := m.migrateMoveDownloadsToStorage(ctx)
	if err != nil {
		hasError = true
	}

	if hasError {
		return fmt.Errorf("errors in the 'downloader usecase' migrations")
	}

	return nil
}
