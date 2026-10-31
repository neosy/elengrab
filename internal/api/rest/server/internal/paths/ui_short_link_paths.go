package httppaths

import (
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	// Paths short links, e.g. /s/{shortCode}
	ShortLinkPath = "/{shortCode}"

	ShortLinkStreamPath        = ShortLinkPath + "/stream"
	ShortLinkWatchTrackingPath = ShortLinkPath + "/watch-tracking"
	ShortLinkImagePath         = ShortLinkPath + "/image"
)

func buildShortLinkPath(shortLinkPrefix, path, shortCode string) string {
	return shortLinkPrefix + strings.Replace(path, "{shortCode}", shortCode, 1)
}

func BuildShortLinkPath(shortLinkPrefix, shortCode string) string {
	return shortLinkPrefix + strings.Replace(ShortLinkPath, "{shortCode}", shortCode, 1)
}

func BuildShortLinkStreamPath(shortLinkPrefix, shortCode string) string {
	return buildShortLinkPath(shortLinkPrefix, ShortLinkStreamPath, shortCode)
}

func BuildShortLinkImagePath(shortLinkPrefix, shortCode string, verHash string, sources []dtypes.ImageSource) string {
	urlSufix := buildImageSufix(verHash, sources)
	return buildShortLinkPath(shortLinkPrefix, ShortLinkImagePath, shortCode) + urlSufix
}
