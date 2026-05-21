package humanize

import (
	"fmt"
)

// DurationClock formats duration in seconds into clock-style representation.
//
// Examples:
//
//	51   -> "0:51"
//	111  -> "1:51"
//	3661 -> "1:01:01"
func DurationClock(durationSec int64) string {
	hr := durationSec / 3600
	min := (durationSec % 3600) / 60
	sec := durationSec % 60

	if hr > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hr, min, sec)
	}

	return fmt.Sprintf("%d:%02d", min, sec)
}
