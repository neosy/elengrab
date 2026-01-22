package ddownload

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type MediaInfo struct {
	FormatType dtypes.FormatType `json:"formatType"`
	Format     dtypes.FileFormat `json:"format"`

	VideoInfo *VideoInfo `json:"videoInfo,omitempty"`
	AudioInfo *AudioInfo `json:"audioInfo,omitempty"`
}

type VideoInfo struct {
	Codec dtypes.VideoCodec `json:"codec"`
	// Bitrate, kbps
	Bitrate    int                    `json:"bitrate"`
	Resolution dtypes.VideoResolution `json:"resolution"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
}

type AudioInfo struct {
	Codec dtypes.AudioCodec `json:"codec"`
	// Bitrate, kbps
	Bitrate int `json:"bitrate"`
	// Hz
	SampleRate *int `json:"sampleRate,omitempty"`
}

func (m *MediaInfo) Copy() *MediaInfo {
	if m == nil {
		return nil
	}

	info := *m
	info.VideoInfo = uptr.Copy(m.VideoInfo)
	info.AudioInfo = uptr.Copy(m.AudioInfo)

	return &info
}
