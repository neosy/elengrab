package idto

import (
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Format struct {
	FormatID   string  `json:"format_id"`
	FileExt    string  `json:"ext"`
	Height     int     `json:"height"`
	Width      int     `json:"width"`
	FPS        float32 `json:"fps"`
	Format     string  `json:"format"`
	FormatNote string  `json:"format_note"`
	Resolution string  `json:"resolution"`
	VCodec     string  `json:"vcodec"`
	Vbr        float32 `json:"vbr"`
	ACodec     string  `json:"acodec"`
	Abr        float32 `json:"abr"`
	Asr        *int    `json:"asr"`
	Filesize   *int    `json:"filesize"`
}

type YouTubeInfo struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	ChannelID  string   `json:"channel_id"`
	ChannelUrl string   `json:"channel_url"`
	Formats    []Format `json:"formats"`
}

func (f *Format) AudioCodec() dtypes.AudioCodec {
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

func (f *Format) VideoCodec() dtypes.VideoCodec {
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
