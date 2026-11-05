package dto

import (
	"fmt"
	"time"
)

type Playlist struct {
	Items []PlaylistItem
}

type PlaylistItem struct {
	ID    string
	Index int

	URL   string
	Title string

	Duration time.Duration
}

// DurationText returns the duration of the playlist item formatted as a string in the format "HH:MM:SS"
// or "MM:SS" if the duration is less than an hour.
func (p *PlaylistItem) DurationText() string {
	var duration time.Duration

	if p != nil {
		duration = p.Duration
	}

	seconds := int64(duration.Seconds())

	hours := seconds / 3600
	minutes := seconds % 3600 / 60
	seconds %= 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
