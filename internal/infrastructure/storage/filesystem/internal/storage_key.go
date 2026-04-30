package core

import "fmt"

func BuildStorageKeyPath(storageKey string, format string) string {
	key := storageKey

	return fmt.Sprintf(
		"%s/%s/%s.%s",
		key[:2],
		key[2:4],
		key,
		format,
	)
}
