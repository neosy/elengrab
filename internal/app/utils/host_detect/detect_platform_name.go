package hostdetect

import dtypes "github.com/neosy/elengrab/internal/domain/types"

func DetectPlatformName(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	platformType := DetectPlatformType(rawURL)
	if platformType != dtypes.MediaPlatformTypeNone {
		return platformType.String()
	}

	return ExtractPlatformName(rawURL)
}
