package fasthttpx

import (
	"github.com/neosy/elengrab/internal/pkg/filex"
)

// SanitizeFileName normalizes a file name to be safe for filesystem storage and HTTP sendfile usage
func SanitizeFileName(name string) string {
	return filex.SanitizeHTTPName(name)
}
