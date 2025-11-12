package ddownload

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// DownloadOptions defines optional parameters for a download operation.
type DownloadOptions struct {
	// Type of content to download (video, audio, or both)
	FormatType dtypes.FormatType `json:"format_type"`

	// Video format (orig, mp4)
	VideoFormat *dtypes.VideoFormat `json:"video_format,omitempty"`

	// Audio format (orig, mp3)
	AudioFormat *dtypes.AudioFormat `json:"audio_format,omitempty"`

	// Custom output directory
	DownloadsDir *string `json:"downloads_dir,omitempty"`

	// Custom file name for the downloaded content
	Filename *string `json:"filename,omitempty"`

	// Desired video quality
	VideoQuality *string `json:"video_quality,omitempty"`

	// Desired audio quality
	AudioQuality *string `json:"audio_quality,omitempty"`
}
