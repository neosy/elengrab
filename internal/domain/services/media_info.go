package dservices

import (
	"strconv"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaInfo struct {
	FormatType dtypes.FormatType
	Format     dtypes.FileFormat

	Duration time.Duration

	VideoInfo *dtypes.VideoInfo
	AudioInfo *dtypes.AudioInfo
}

func NewMediaInfoFromDomain(dMediaInfo *dtypes.MediaInfo) *MediaInfo {
	if dMediaInfo == nil {
		return nil
	}

	duration, _ := strconv.ParseFloat(dMediaInfo.Duration, 64)

	return &MediaInfo{
		FormatType: dMediaInfo.FormatType,
		Format:     dMediaInfo.Format,

		Duration: time.Duration(duration * float64(time.Second)),

		VideoInfo: dMediaInfo.VideoInfo,
		AudioInfo: dMediaInfo.AudioInfo,
	}
}

func (m MediaInfo) MediaInfoDomain() dtypes.MediaInfo {
	var duration string
	if m.Duration > 0 {
		duration = strconv.FormatFloat(m.Duration.Seconds(), 'f', 6, 64)
	}

	return dtypes.MediaInfo{
		FormatType: m.FormatType,
		Format:     m.Format,

		Duration:   duration,
		DurationMs: m.Duration.Milliseconds(),

		VideoInfo: uptr.Copy(m.VideoInfo),
		AudioInfo: uptr.Copy(m.AudioInfo),
	}
}

func (m *MediaInfo) MediaInfoDomainPtr() *dtypes.MediaInfo {
	if m == nil {
		return nil
	}
	return new(m.MediaInfoDomain())
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

// HasVideo
func (m *MediaInfo) HasVideo() bool {
	return m.VideoInfo != nil
}
