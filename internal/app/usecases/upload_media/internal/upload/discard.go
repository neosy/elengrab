package upload

import (
	"errors"
	"io/fs"
	"os"
)

// Discard removes the temporary uploaded files. It is safe to call after any outcome:
// once the file has been moved to the storage, there is nothing left to remove.
func (u *Upload) Discard() {
	for _, path := range []string{u.PartPath(), u.FilePath()} {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			u.logger.Warn("Failed to remove temporary uploaded file", "path", path, "error", err)
		}
	}
}
