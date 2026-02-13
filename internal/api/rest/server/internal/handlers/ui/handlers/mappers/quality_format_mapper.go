package mappers

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

var (
	videoFormatMap = map[string]dtypes.VideoFormat{
		"auto": dtypes.VideoFormatAuto,
		"mp4":  dtypes.VideoFormatMP4,
		"webm": dtypes.VideoFormatWebM,
	}
	videoCodecMap = map[string]dtypes.VideoCodec{
		"best": dtypes.VideoCodecBest,
		"h264": dtypes.VideoCodecH264,
		"h265": dtypes.VideoCodecH265,
		"av1":  dtypes.VideoCodecAV1,
	}
	videoCodecWebMMap = map[string]struct{}{
		"best": {},
		"av1":  {},
	}
	videoResolutionMap = map[string]dtypes.VideoResolution{
		"max":  dtypes.VideoResolutionMax,
		"1080": dtypes.VideoResolution1080p,
		"720":  dtypes.VideoResolution720p,
		"480":  dtypes.VideoResolution480p,
		"360":  dtypes.VideoResolution360p,
	}
	audioFormatMap = map[string]dtypes.AudioFormat{
		"auto": dtypes.AudioFormatAuto,
		"mp3":  dtypes.AudioFormatMP3,
		"m4a":  dtypes.AudioFormatM4A,
		"flac": dtypes.AudioFormatFLAC,
		"opus": dtypes.AudioFormatOPUS,
	}
	onlyAudioFormatMap = map[string]dtypes.AudioFormat{
		"mp3":  dtypes.AudioFormatMP3,
		"m4a":  dtypes.AudioFormatM4A,
		"flac": dtypes.AudioFormatFLAC,
		"opus": dtypes.AudioFormatOPUS,
	}
)

func (m *Mappers) MapFormatType(qc, f string) dtypes.FormatType {
	switch qc {
	case "best":
		if _, exists := onlyAudioFormatMap[f]; exists {
			return dtypes.FormatTypeAudioOnly
		}
	case "only_audio":
		return dtypes.FormatTypeAudioOnly
	}

	return dtypes.FormatTypeVideoAudio
}

func (m *Mappers) MapVideoFormat(qc, f string) *dtypes.VideoFormat {
	var videoFormat *dtypes.VideoFormat
	if _, exists := videoFormatMap[f]; exists {
		videoFormat = uptr.Any(videoFormatMap[f])
	}
	if videoFormat != nil && *videoFormat == dtypes.VideoFormatWebM {
		if _, exists := videoCodecWebMMap[qc]; !exists {
			videoFormat = dtypes.VideoFormatAuto.Ptr()
		}
	}
	return videoFormat
}

func (m *Mappers) MapVideoCodec(qc string) *dtypes.VideoCodec {
	var videoCodec *dtypes.VideoCodec
	if codec, exists := videoCodecMap[qc]; exists {
		videoCodec = &codec
	}
	return videoCodec
}

func (m *Mappers) MapVideoResolution(qr string) *dtypes.VideoResolution {
	var videoResolution *dtypes.VideoResolution
	if res, exists := videoResolutionMap[qr]; exists {
		videoResolution = uptr.Any(res)
	}
	return videoResolution
}

func (m *Mappers) MapAudioFormat(f string) *dtypes.AudioFormat {
	var audioFormat *dtypes.AudioFormat
	if _, exists := audioFormatMap[f]; exists {
		audioFormat = uptr.Any(audioFormatMap[f])
	}
	return audioFormat
}
