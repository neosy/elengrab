package core

import (
	"path/filepath"
)

func BuildStorageKeyPath(storageKey string, format string) string {
	key := storageKey

	return filepath.Join(
		key[:2],
		key[2:4],
		key+"."+format,
	)
}
