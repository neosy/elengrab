package dto

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type DLOptions struct {
	FormatType             dtypes.FormatType
	VideoFormat            dtypes.VideoFormat
	VideoCodec             dtypes.VideoCodec
	VideoResolution        dtypes.VideoResolution
	AudioFormat            dtypes.AudioFormat
	ConcurrentFragments    uint8
	RequiresYouTubeCookies bool
}
