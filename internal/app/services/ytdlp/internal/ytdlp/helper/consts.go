package helper

import dtypes "github.com/neosy/elengrab/internal/domain/types"

const (
	formatTypeDefault = dtypes.FormatTypeVideoAudio

	videoFormatDefault     = dtypes.VideoFormatMP4
	videoCodecDefault      = dtypes.VideoCodecBest
	videoResolutionDefault = dtypes.VideoResolutionBest

	audioFormatDefault      = dtypes.AudioFormatMP3
	audioQualityMP3Default  = "2"
	audioQualityM4ADefault  = "2"
	audioQualityOPUSDefault = "160K"

	concurrentFragmentsDefault = 5
)
