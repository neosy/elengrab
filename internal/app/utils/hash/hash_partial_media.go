package hash

import (
	"path/filepath"

	nfile "github.com/neosy/elengrab/internal/pkg/filex"
)

const (
	partialHashBlockSize = 100 * 1024
	partialHashBlocks    = 3
)

// FilePartialHash selects a partial-hash method for media files.
// webm/opus → skip first 1 KB; others → standard partial hash.
func FilePartialHash(filePath string) (string, error) {
	ext := filepath.Ext(filePath)

	switch ext[1:] {
	case "webm", "opus":
		return nfile.PartialHashWithOffset(filePath, 1, 1024, 1024)
	default:
		return nfile.PartialHash(filePath, partialHashBlocks, partialHashBlockSize)
	}
}
