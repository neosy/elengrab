package dservices

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// DownloadOptions defines optional parameters for a download operation.
type DownloadOptions struct {
	// Type of content to download (video, audio, or both)
	FormatType dtypes.FormatType

	// Video format (orig, mp4)
	VideoFormat *dtypes.VideoFormat

	// Video codec (best, h264, ...)
	VideoCodec *dtypes.VideoCodec

	// Video resolution (best, 4k, 2k, 1080p, 720p, ...)
	VideoResolution *dtypes.VideoResolution

	// Audio format (orig, mp3)
	AudioFormat *dtypes.AudioFormat

	// Custom output directory
	DownloadsDir *string

	// Custom file name for the downloaded content
	Filename *string

	// IncludeTitleInFilename indicates whether to include the title in the file name
	IncludeTitleInFilename bool

	// Desired video quality
	VideoQuality *string

	// Desired audio quality
	AudioQuality *string
}
