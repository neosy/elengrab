package dtypes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type MediaInfo struct {
	Format     FileFormat `json:"format"`
	FormatType FormatType `json:"formatType"`

	DurationText string `json:"duration,omitempty"`
	DurationMs   int64  `json:"durationMs,omitempty"`

	Bitrate int `json:"bitrate,omitempty"`

	VideoInfo *VideoInfo `json:"videoInfo,omitempty"`
	AudioInfo *AudioInfo `json:"audioInfo,omitempty"`

	ThumbnailID      *uuid.UUID `json:"thumbnailId,omitempty"`
	FrameThumbnailID *uuid.UUID `json:"frameThumbnailId,omitempty"`
}

func NewMediaInfo(ext string) *MediaInfo {
	format := MapFileExtToFileFormat(ext)

	return &MediaInfo{
		Format:     format,
		FormatType: format.FormatType(),
	}
}

func (info *MediaInfo) SetDuration(duration time.Duration) {
	if info == nil {
		return
	}

	info.DurationText = strconv.FormatFloat(duration.Seconds(), 'f', 6, 64)
	info.DurationMs = int64(duration.Milliseconds())
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

// Duration returns the duration of the media as a time.Duration.
func (m *MediaInfo) Duration() time.Duration {
	return time.Duration(m.DurationMs) * time.Millisecond
}

// IsShorts returns true if the media is a video and its duration is less than or equal to 3 minutes (180 seconds).
func (m *MediaInfo) IsShorts() bool {
	if m == nil {
		return false
	}

	if m.VideoInfo == nil {
		return false
	}

	if !m.IsVideo() {
		return false
	}

	duration := m.Duration()

	isShorts := duration > 0 &&
		duration <= 180*time.Second &&
		m.VideoInfo.AspectRatio() <= 1.0

	return isShorts
}

// IsAudioOnly returns true if the media is audio-only (no video).
func (m *MediaInfo) IsAudioOnly() bool {
	return m != nil && m.FormatType == FormatTypeAudioOnly
}

// IsVideo returns true if the media has video (not audio-only).
func (m *MediaInfo) IsVideo() bool {
	return m != nil && m.VideoInfo != nil
}

// ResolutionString returns the video dimensions formatted as "widthxheight" (e.g., "1920x1080").
func (v *VideoInfo) ResolutionString() string {
	return fmt.Sprintf("%dx%d", v.Width, v.Height)
}

// AspectRatio returns the aspect ratio of the video as a float64 (width / height).
func (v *VideoInfo) AspectRatio() float64 {
	if v == nil || v.Width == 0 || v.Height == 0 {
		return 0
	}
	return float64(v.Width) / float64(v.Height)
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

func (a *AudioInfo) Merge(src AudioInfo) bool {
	var updated bool

	if src.Codec != "" && a.Codec != src.Codec {
		a.Codec = src.Codec
		updated = true
	}

	if src.Bitrate != 0 && a.Bitrate != src.Bitrate {
		a.Bitrate = src.Bitrate
		updated = true
	}

	if src.SampleRate != nil {
		if a.SampleRate == nil || *a.SampleRate != *src.SampleRate {
			a.SampleRate = src.SampleRate
			updated = true
		}
	}

	return updated
}

func (v *VideoInfo) Merge(src VideoInfo) bool {
	var updated bool

	if src.Codec != "" && v.Codec != src.Codec {
		v.Codec = src.Codec
		updated = true
	}

	if src.Bitrate != 0 && v.Bitrate != src.Bitrate {
		v.Bitrate = src.Bitrate
		updated = true
	}

	if src.Resolution != "" && v.Resolution != src.Resolution {
		v.Resolution = src.Resolution
		updated = true
	}

	if src.Width != 0 && v.Width != src.Width {
		v.Width = src.Width
		updated = true
	}

	if src.Height != 0 && v.Height != src.Height {
		v.Height = src.Height
		updated = true
	}

	return updated
}

func (vi *VideoInfo) IsPortrait() bool {
	if vi == nil {
		return false
	}

	return vi.Width < vi.Height
}
