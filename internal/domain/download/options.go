package ddownload

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// DownloadOptions defines optional parameters for a download operation.
type DownloadOptions struct {
	// Type of content to download (video, audio, or both)
	FormatType dtypes.FormatType

	// Custom output directory
	DownloadsDir *string

	// Custom file name for the downloaded content
	Filename *string

	// Desired video quality
	VideoQuality *string

	// Desired audio quality
	AudioQuality *string
}
