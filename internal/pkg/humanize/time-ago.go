package humanize

import (
	"fmt"
	"time"
)

const (
	timeUnitSecond = "second"
	timeUnitMinute = "minute"
	timeUnitHour   = "hour"
	timeUnitDay    = "day"
	timeUnitWeek   = "week"
	timeUnitMonth  = "month"
	timeUnitYear   = "year"
)

func TimeAgo(t time.Time) string {
	var (
		timeUnit, prefix string
		timeInterval     int
	)

	now := time.Now()
	if t.After(now) {
		t = now
	}

	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		timeInterval = int(diff / time.Second)
		timeUnit = timeUnitSecond
	case diff < time.Hour:
		timeInterval = int(diff / time.Minute)
		timeUnit = timeUnitMinute
	case diff < 24*time.Hour:
		timeInterval = int(diff / time.Hour)
		timeUnit = timeUnitHour
	case diff < 7*24*time.Hour:
		timeInterval = int(diff / (24 * time.Hour))
		timeUnit = timeUnitDay
	case diff < 30*24*time.Hour:
		timeInterval = int(diff / (7 * 24 * time.Hour))
		timeUnit = timeUnitWeek
	case diff < 365*24*time.Hour:
		months := 0
		tmp := t
		for tmp.AddDate(0, 1, 0).Before(now) || tmp.AddDate(0, 1, 0).Equal(now) {
			tmp = tmp.AddDate(0, 1, 0)
			months++
		}
		timeInterval = months
		timeUnit = timeUnitMonth
	default:
		years := 0
		tmp := t
		for tmp.AddDate(1, 0, 0).Before(now) || tmp.AddDate(1, 0, 0).Equal(now) {
			tmp = tmp.AddDate(1, 0, 0)
			years++
		}
		timeInterval = years
		timeUnit = timeUnitYear
	}

	if timeInterval <= 0 {
		timeInterval = 1
		prefix = "less "
	}

	if timeInterval > 1 {
		timeUnit += "s"
	}

	suffix := " ago"

	return fmt.Sprintf("%s%d %s%s", prefix, timeInterval, timeUnit, suffix)
}
