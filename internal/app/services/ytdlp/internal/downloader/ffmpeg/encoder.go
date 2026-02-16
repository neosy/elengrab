package ffmpeg

import dtypes "github.com/neosy/elengrab/internal/domain/types"

const (
	defaultVideoH264 = "libx264 -crf 22 -preset slow"
	defaultVideoH265 = "libx265 -crf 22 -preset slow"
	defaultVideoAV1  = "libaom-av1 -crf 0 -b:v 0"
)

func VideoEncoderArgs(videoCode dtypes.VideoCodec) string {
	switch videoCode {
	case dtypes.VideoCodecH264:
		return defaultVideoH264
	case dtypes.VideoCodecAV1:
		return defaultVideoAV1
	case dtypes.VideoCodecH265:
		return defaultVideoH265
	default:
		return defaultVideoH264
	}
}
