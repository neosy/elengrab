package utils

import (
	"path/filepath"

	"github.com/neosy/elengrab/internal/pkg/nfile"
)

const (
	hashPartialBlockSize = 100 * 1024
	hashPartialBlocks    = 3
)

// HashPartialMedia selects a partial-hash method for media files.
// webm/opus → skip first 1 KB; others → standard partial hash.
func HashPartialMedia(filePath string) (string, error) {
	ext := filepath.Ext(filePath)

	switch ext[1:] {
	case "webm", "opus":
		return nfile.HashPartialWithOffset(filePath, 1, 1024, 1024)
	default:
		return nfile.HashPartial(filePath, hashPartialBlocks, hashPartialBlockSize)
	}
}
