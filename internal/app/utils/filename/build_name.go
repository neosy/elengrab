package filename

import (
	"strings"

	"uuid"

	"github.com/neosy/elengrab/internal/pkg/filex"
	"github.com/neosy/elengrab/internal/pkg/stringx"
)

const (
	MaxTitleLength = 100
)

// BuildUnique builds a unique file name: "<title>_<uniqueName>.<ext>",
// or "<uniqueName>.<ext>" when title is empty. If uniqueName is empty, a UUID is generated.
// The title is sanitized and truncated to fit; uniqueName and ext are used as is,
// so the result stays within filex.NameMaxBytes only if they fit on their own.
func BuildUnique(uniqueName, title, ext string) string {
	if uniqueName == "" {
		uniqueName = uuid.New().String()
	}

	suffix := ext
	if !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}

	name := uniqueName

	if title != "" {
		safeTitle := filex.SanitizeHTTPName(title)

		// Bytes left for the title after the unique name, "_" and the extension.
		maxTitleBytes := min(
			MaxTitleLength,
			filex.NameMaxBytes-len(uniqueName)-len("_")-len(suffix),
		)

		if maxTitleBytes > 0 {
			safeTitle = stringx.TruncateBytesWords(safeTitle, maxTitleBytes)
			if safeTitle != "" {
				name = safeTitle + "_" + uniqueName
			}
		}
	}

	return name + suffix
}
