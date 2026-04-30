package ffmpegsrv

import dtypes "github.com/neosy/elengrab/internal/domain/types"

// FrameFormat defines the interface for different frame output formats.
type FrameFormat interface {
	// Args returns the ffmpeg arguments for the frame output format.
	Args() []string
	// Format returns the image format type for the frame output.
	Format() dtypes.ImageFormat
}

// FrameFormatJPEG defines the JPEG output format for frames.
type FrameFormatJPEG struct{}

func (f FrameFormatJPEG) Args() []string {
	return []string{
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
	}
}

func (f FrameFormatJPEG) Format() dtypes.ImageFormat {
	return dtypes.ImageFormatJPEG
}

// FrameFormatPNG defines the PNG output format for frames.
type FrameFormatPNG struct{}

func (f FrameFormatPNG) Args() []string {
	return []string{
		"-f", "image2pipe",
		"-vcodec", "png",
	}
}

func (f FrameFormatPNG) Format() dtypes.ImageFormat {
	return dtypes.ImageFormatPNG
}

// FrameFormatWebP defines the WebP output format for frames.
type FrameFormatWebP struct{}

func (f FrameFormatWebP) Args() []string {
	return []string{
		"-f", "image2pipe",
		"-vcodec", "libwebp",
	}
}

func (f FrameFormatWebP) Format() dtypes.ImageFormat {
	return dtypes.ImageFormatWebP
}
