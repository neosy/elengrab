package backupname

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/neosy/elengrab/internal/pkg/filex"
)

// Prefix returns the sanitized lowercase file name prefix for the given database.
func Prefix(dbName string) string {
	prefix := filex.SanitizeFileName(dbName)
	prefix = strings.ToLower(prefix)
	return prefix
}

// Path builds the full backup file path for the given database.
func Path(dbName, backupsDir string) string {
	filename := fmt.Sprintf(
		"%s_%s.%s",
		Prefix(dbName),
		time.Now().Format(DateFormat),
		FileExt,
	)

	return filepath.Join(backupsDir, filename)
}
