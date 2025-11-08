package dyoutubeinfo

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type Format struct {
	// Type of the format (e.g., only_audio, only_video, video_audio)
	FormatType dtypes.FormatType

	FormatId string

	// File extension (e.g., mp4, webm, m4a)
	FileExt string

	// Video height in pixels
	Height int

	// Video width in pixels
	Width int

	// Frames per second (optional)
	FPS *int

	// Full format description
	Format string

	// Additional format note (e.g., "medium", "high", "storyboard")
	FormatNote string

	// Resolution string (e.g., "1920x1080")
	Resolution string

	// Video codec (optional)
	VCodec *string

	// Audio codec (optional)
	ACodec *string

	// Video bitrate in kbit/s (optional)
	Vbr *float32

	// Audio bitrate in kbit/s (optional)
	Abr *float32

	// Audio sample rate in Hz (optional)
	Asr *int

	// File size in bytes
	Filesize *int
}
