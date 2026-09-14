package ddownload

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// DownloadOptions defines optional parameters for a download operation.
type DownloadOptions struct {
	// Type of content to download (video, audio, or both)
	FormatType dtypes.FormatType `json:"format_type"`

	// Video format (best, mp4)
	VideoFormat *dtypes.VideoFormat `json:"video_format,omitempty"`

	// Video codec (best, h264, ...)
	VideoCodec *dtypes.VideoCodec `json:"video_codec,omitempty"`

	// Video resolution (best, 4k, 2k, 1080p, 720p, ...)
	VideoResolution *dtypes.VideoResolution

	// Audio format (orig, mp3)
	AudioFormat *dtypes.AudioFormat `json:"audio_format,omitempty"`

	// Custom file name for the downloaded content
	Filename *string `json:"filename,omitempty"`

	// Desired video quality
	VideoQuality *string `json:"video_quality,omitempty"`

	// Desired audio quality
	AudioQuality *string `json:"audio_quality,omitempty"`
}

func (src *DownloadOptions) Clone() *DownloadOptions {
	if src == nil {
		return nil
	}

	copy := *src

	copy.VideoFormat = uptr.Clone(src.VideoFormat)
	copy.VideoCodec = uptr.Clone(src.VideoCodec)
	copy.VideoResolution = uptr.Clone(src.VideoResolution)
	copy.AudioFormat = uptr.Clone(src.AudioFormat)
	copy.Filename = uptr.Clone(src.Filename)
	copy.VideoQuality = uptr.Clone(src.VideoQuality)
	copy.AudioQuality = uptr.Clone(src.AudioQuality)

	return &copy
}
