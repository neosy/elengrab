package pservices

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// FrameFormat defines the interface for different frame output formats.
type FrameFormat interface {
	// Args returns the ffmpeg arguments for the frame output format.
	Args() []string
	// Format returns the image format type for the frame output.
	Format() dtypes.ImageFormat
}

// FrameStrategy defines the interface for different frame extraction strategies.
type FrameStrategy interface {
	// Args returns the ffmpeg arguments for the frame extraction strategy.
	Args() []string
}
