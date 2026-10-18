package upload

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Cleanup removes the temporary uploaded file after it has been moved to the storage.
// An already missing file is not an error.
func (u *Upload) Cleanup() error {
	path := u.FilePath()

	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		u.logger.Warn("Failed to remove temporary uploaded file after moving it to the storage", "path", path, "error", err)
		return fmt.Errorf("remove temporary uploaded file %q: %w", path, err)
	}

	return nil
}
