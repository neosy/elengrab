package pservices

import (
	"strconv"
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// MediaInfo is the structured description of a media file:
// duration, bitrate and its video and audio streams.
type MediaInfo struct {
	FormatType dtypes.FormatType
	Format     dtypes.FileFormat

	Duration time.Duration

	// Bitrate, kbps
	Bitrate int

	VideoInfo *dtypes.VideoInfo
	AudioInfo *dtypes.AudioInfo
}

// NewMediaInfoFromExt creates a MediaInfo from a file extension (e.g. "mp4").
func NewMediaInfoFromExt(extNoDot string) *MediaInfo {
	fileFormat := dtypes.FileExtToFileFormat(extNoDot)

	return &MediaInfo{
		FormatType: fileFormat.FormatType(),
		Format:     fileFormat,
	}
}

// NewMediaInfoFromDomain creates a MediaInfo from the domain dtypes.MediaInfo.
func NewMediaInfoFromDomain(dMediaInfo *dtypes.MediaInfo) *MediaInfo {
	if dMediaInfo == nil {
		return nil
	}

	duration, _ := strconv.ParseFloat(dMediaInfo.DurationText, 64)

	return &MediaInfo{
		FormatType: dMediaInfo.FormatType,
		Format:     dMediaInfo.Format,

		Duration: time.Duration(duration * float64(time.Second)),

		VideoInfo: dMediaInfo.VideoInfo,
		AudioInfo: dMediaInfo.AudioInfo,
	}
}

// ToDomain converts the MediaInfo to the domain dtypes.MediaInfo.
func (m *MediaInfo) ToDomain() *dtypes.MediaInfo {
	if m == nil {
		return nil
	}

	var durationText string
	if m.Duration > 0 {
		durationText = strconv.FormatFloat(m.Duration.Seconds(), 'f', 6, 64)
	}

	return &dtypes.MediaInfo{
		Format:     m.Format,
		FormatType: m.FormatType,

		DurationText: durationText,
		DurationMs:   m.Duration.Milliseconds(),

		Bitrate: m.Bitrate,

		VideoInfo: uptr.Clone(m.VideoInfo),
		AudioInfo: uptr.Clone(m.AudioInfo),
	}
}

func (m *MediaInfo) Clone() *MediaInfo {
	if m == nil {
		return nil
	}

	info := *m

	info.VideoInfo = uptr.Clone(m.VideoInfo)
	info.AudioInfo = uptr.Clone(m.AudioInfo)

	return &info
}

// HasVideo
func (m *MediaInfo) HasVideo() bool {
	return m.VideoInfo != nil
}

// DurationSecondsString
func (m *MediaInfo) DurationSecondsString() string {
	if m == nil {
		return ""
	}
	return strconv.FormatFloat(m.Duration.Seconds(), 'f', 6, 64)
}
