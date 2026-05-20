package helper

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	formatTypeDefault = dtypes.FormatTypeVideoAudio

	videoFormatDefault     = dtypes.VideoFormatMP4
	videoCodecDefault      = dtypes.VideoCodecBest
	videoResolutionDefault = dtypes.VideoResolutionMax

	audioFormatDefault = dtypes.AudioFormatMP3

	audioQualityMP3Default         = "2"
	audioQualityMP3BitrateDefault  = 192
	audioQualityM4ADefault         = "2"
	audioQualityAACBitrateDefault  = "160k"
	audioQualityOPUSDefault        = "160k"
	audioQualityFLACBitrateDefault = 965
)
