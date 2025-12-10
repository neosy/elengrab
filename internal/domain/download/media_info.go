package ddownload

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type MediaInfo struct {
	FormatType dtypes.FormatType
	Format     dtypes.FileFormat
	VideoCodec dtypes.VideoCodec
	Resolution dtypes.VideoResolution
	Width      int
	Height     int
}
