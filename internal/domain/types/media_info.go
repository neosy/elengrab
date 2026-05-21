package dtypes

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaInfo struct {
	FormatType FormatType `json:"formatType"`
	Format     FileFormat `json:"format"`

	Duration   string `json:"duration,omitempty"`
	DurationMs int64  `json:"durationMs,omitempty"`

	VideoInfo *VideoInfo `json:"videoInfo,omitempty"`
	AudioInfo *AudioInfo `json:"audioInfo,omitempty"`

	ThumbnailID      *uuid.UUID `json:"thumbnailId,omitempty"`
	FrameThumbnailID *uuid.UUID `json:"frameThumbnailId,omitempty"`
}

type VideoInfo struct {
	Codec VideoCodec `json:"codec"`

	// Bitrate, kbps
	Bitrate int `json:"bitrate"`

	Resolution VideoResolution `json:"resolution"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
}

type AudioInfo struct {
	Codec AudioCodec `json:"codec"`

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
	info.ThumbnailID = uptr.Copy(m.ThumbnailID)
	info.FrameThumbnailID = uptr.Copy(m.FrameThumbnailID)

	return &info
}

// PreferredThumbnailID returns the ThumbnailID if available, otherwise returns the FrameThumbnailID.
func (m *MediaInfo) PreferredThumbnailID() uuid.UUID {
	if m == nil {
		return uuid.Nil
	}

	if m.ThumbnailID != nil {
		return *m.ThumbnailID
	}

	if m.FrameThumbnailID != nil {
		return *m.FrameThumbnailID
	}

	return uuid.Nil
}

// HasVideo
func (m *MediaInfo) HasVideo() bool {
	return m.VideoInfo != nil
}

func (m *MediaInfo) FormatDuration() string {
	return humanize.DurationClock(m.DurationMs / 1000)
}

func (m *MediaInfo) IsPortrait() bool {
	return m.VideoInfo.IsPortrait()
}

// ResolutionString returns the video dimensions formatted as "widthxheight" (e.g., "1920x1080").
func (v *VideoInfo) ResolutionString() string {
	return fmt.Sprintf("%dx%d", v.Width, v.Height)
}

func (vi *VideoInfo) Copy() *VideoInfo {
	if vi == nil {
		return nil
	}
	return new(*vi)
}

func (ai *AudioInfo) Copy() *AudioInfo {
	if ai == nil {
		return nil
	}
	return new(*ai)
}

func (vi *VideoInfo) IsPortrait() bool {
	if vi == nil {
		return false
	}

	return vi.Width < vi.Height
}
