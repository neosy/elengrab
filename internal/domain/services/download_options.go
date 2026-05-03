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

	// Custom file name for the downloaded content
	Filename *string

	// IncludeTitleInFilename indicates whether to include the title in the file name
	IncludeTitleInFilename bool

	// Desired video quality
	VideoQuality *string

	// Desired audio quality
	AudioQuality *string

	// DownloadChannelAvatar indicates whether to download the channel's avatar
	DownloadChannelAvatar bool
}

// NewDefaultDownloadOptions creates a DownloadOptions struct with default values.
func NewDefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		IncludeTitleInFilename: true,
		DownloadChannelAvatar:  true,
	}
}
