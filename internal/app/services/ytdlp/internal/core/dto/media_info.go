package idto

import (
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaFormat struct {
	FormatID       string  `json:"format_id"`
	FileExt        string  `json:"ext"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	FPS            float32 `json:"fps"`
	Format         string  `json:"format"`
	FormatNote     string  `json:"format_note"`
	Resolution     string  `json:"resolution"`
	VCodec         string  `json:"vcodec"`
	Vbr            float32 `json:"vbr"`
	ACodec         string  `json:"acodec"`
	Abr            float32 `json:"abr"`
	Asr            *int    `json:"asr"`
	Tbr            float32 `json:"tbr"`
	Filesize       *int64  `json:"filesize"`
	FilesizeApprox *int64  `json:"filesize_approx"`
}

type MediaInfo struct {
	ID         string        `json:"id"`
	Title      string        `json:"title"`
	Extractor  string        `json:"extractor"`
	ChannelID  string        `json:"channel_id"`
	ChannelUrl string        `json:"channel_url"`
	Duration   int           `json:"duration"`
	Formats    []MediaFormat `json:"formats"`
}

func (f *MediaFormat) AudioCodec() dtypes.AudioCodec {
	if f.ACodec == "" {
		return dtypes.AudioCodecNone
	}
	if strings.HasPrefix(f.ACodec, "mp4a") {
		return dtypes.AudioCodecAAC
	}
	if strings.HasPrefix(f.ACodec, "opus") {
		return dtypes.AudioCodecOPUS
	}
	return dtypes.AudioCodecNone
}

func (f *MediaFormat) VideoCodec() dtypes.VideoCodec {
	if f.VCodec == "" {
		return dtypes.VideoCodecNone
	}
	if strings.HasPrefix(f.VCodec, "av01") {
		return dtypes.VideoCodecAV1
	}
	if strings.HasPrefix(f.VCodec, "vp9") {
		return dtypes.VideoCodecVP9
	}
	if strings.HasPrefix(f.VCodec, "avc1") {
		return dtypes.VideoCodecH264
	}
	return dtypes.VideoCodecNone
}
