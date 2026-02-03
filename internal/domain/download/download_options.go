package ddownload

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
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

func (src *DownloadOptions) Copy() *DownloadOptions {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)
	copy.VideoFormat = uptr.Copy(src.VideoFormat)
	copy.VideoCodec = uptr.Copy(src.VideoCodec)
	copy.VideoResolution = uptr.Copy(src.VideoResolution)
	copy.AudioFormat = uptr.Copy(src.AudioFormat)
	copy.Filename = uptr.Copy(src.Filename)
	copy.VideoQuality = uptr.Copy(src.VideoQuality)
	copy.AudioQuality = uptr.Copy(src.AudioQuality)

	return copy
}
