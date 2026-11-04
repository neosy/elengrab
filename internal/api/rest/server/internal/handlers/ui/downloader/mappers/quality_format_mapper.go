package mappers

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var (
	videoFormatByStringValue = map[string]dtypes.VideoFormat{
		"auto": dtypes.VideoFormatAuto,
		"mp4":  dtypes.VideoFormatMP4,
		"webm": dtypes.VideoFormatWebM,
	}
	videoCodecByStringValue = map[string]dtypes.VideoCodec{
		"best": dtypes.VideoCodecBest,
		"h264": dtypes.VideoCodecH264,
		"h265": dtypes.VideoCodecH265,
		"av1":  dtypes.VideoCodecAV1,
	}
	videoCodecWebMAllowed = map[string]struct{}{
		"best": {},
		"av1":  {},
	}
	videoResolutionByStringValue = map[string]dtypes.VideoResolution{
		"max":  dtypes.VideoResolutionMax,
		"1080": dtypes.VideoResolution1080p,
		"720":  dtypes.VideoResolution720p,
		"480":  dtypes.VideoResolution480p,
		"360":  dtypes.VideoResolution360p,
	}
	audioFormatByStringValue = map[string]dtypes.AudioFormat{
		"auto": dtypes.AudioFormatAuto,
		"mp3":  dtypes.AudioFormatMP3,
		"m4a":  dtypes.AudioFormatM4A,
		"flac": dtypes.AudioFormatFLAC,
		"opus": dtypes.AudioFormatOPUS,
	}
	onlyAudioFormatByStringValue = map[string]dtypes.AudioFormat{
		"mp3":  dtypes.AudioFormatMP3,
		"m4a":  dtypes.AudioFormatM4A,
		"flac": dtypes.AudioFormatFLAC,
		"opus": dtypes.AudioFormatOPUS,
	}
)

func (m *Mappers) MapFormatType(qc, f string) dtypes.FormatType {
	switch qc {
	case "best":
		if _, exists := onlyAudioFormatByStringValue[f]; exists {
			return dtypes.FormatTypeAudioOnly
		}
	case "only_audio":
		return dtypes.FormatTypeAudioOnly
	}

	return dtypes.FormatTypeVideoAudio
}

func (m *Mappers) MapVideoFormat(qc, f string) dtypes.VideoFormat {
	videoFormat := dtypes.VideoFormatNone

	if _, exists := videoFormatByStringValue[f]; exists {
		videoFormat = videoFormatByStringValue[f]
	}
	
	if videoFormat == dtypes.VideoFormatWebM {
		if _, exists := videoCodecWebMAllowed[qc]; !exists {
			videoFormat = dtypes.VideoFormatAuto
		}
	}
	return videoFormat
}

func (m *Mappers) MapVideoCodec(qc string) dtypes.VideoCodec {
	videoCode, exists := videoCodecByStringValue[qc]
	if !exists {
		return dtypes.VideoCodecNone
	}

	return videoCode
}

func (m *Mappers) MapVideoResolution(qr string) dtypes.VideoResolution {
	res, exists := videoResolutionByStringValue[qr]
	if !exists {
		return dtypes.VideoResolutionNone
	}

	return res
}

func (m *Mappers) MapAudioFormat(f string) dtypes.AudioFormat {
	audioFormat, exists := audioFormatByStringValue[f]
	if !exists {
		return dtypes.AudioFormatNone
	}

	return audioFormat
}
