package uformat

import (
	"fmt"
	"time"
)

type DurationUnit uint8

const (
	DurationUnitNone DurationUnit = iota
	DurationUnitNanosecond
	DurationUnitMicrosecond
	DurationUnitMillisecond
	DurationUnitSecond
	DurationUnitMinute
	DurationUnitHour
)

var (
	mapDurationUnitToShortName = map[DurationUnit]string{
		DurationUnitNanosecond:  "ns",
		DurationUnitMicrosecond: "µs",
		DurationUnitMillisecond: "ms",
		DurationUnitSecond:      "s",
		DurationUnitMinute:      "m",
		DurationUnitHour:        "h",
	}

	mapDurationUnitShortName = map[string]DurationUnit{
		"ns": DurationUnitNanosecond,
		"µs": DurationUnitMicrosecond,
		"ms": DurationUnitMillisecond,
		"s":  DurationUnitSecond,
		"m":  DurationUnitMinute,
		"h":  DurationUnitHour,
	}

	mapDurationUnitToFullName = map[DurationUnit]string{
		DurationUnitNanosecond:  "nanoseconds",
		DurationUnitMicrosecond: "microseconds",
		DurationUnitMillisecond: "milliseconds",
		DurationUnitSecond:      "seconds",
		DurationUnitMinute:      "minutes",
		DurationUnitHour:        "hours",
	}
)

// DurationUnitByValue returns the duration unit type based on the given duration.
func DurationUnitByValue(d time.Duration) DurationUnit {
	abs := d.Abs()

	switch {
	case abs < time.Microsecond:
		return DurationUnitNanosecond
	case abs < time.Millisecond:
		return DurationUnitMicrosecond
	case abs < time.Second:
		return DurationUnitMillisecond
	case abs < time.Minute:
		return DurationUnitSecond
	case abs < time.Hour:
		return DurationUnitMinute
	case abs < 24*time.Hour:
		return DurationUnitHour
	default:
		return DurationUnitNone
	}
}

// ToFloat64 converts a time.Duration to a float64 based on the given duration unit.
func (u DurationUnit) ToFloat64(d time.Duration) float64 {
	switch u {
	case DurationUnitNanosecond:
		return float64(d)
	case DurationUnitMicrosecond:
		return float64(d) / float64(time.Microsecond)
	case DurationUnitMillisecond:
		return float64(d) / float64(time.Millisecond)
	case DurationUnitSecond:
		return float64(d) / float64(time.Second)
	case DurationUnitMinute:
		return float64(d) / float64(time.Minute)
	case DurationUnitHour:
		return float64(d) / float64(time.Hour)
	default:
		return float64(d)
	}
}

// Unit returns the string representation of the duration unit type.
func (u DurationUnit) String() string {
	return u.Name()
}

// Name returns the short name of the duration unit type.
func (u DurationUnit) Name() string {
	return mapDurationUnitToShortName[u]
}

// FullName returns the full name of the duration unit type.
func (u DurationUnit) FullName() string {
	return mapDurationUnitToFullName[u]
}

// ParseDurationUnit parses a string representation of a duration unit into the corresponding DurationUnit type.
func ParseDurationUnit(unit string) (DurationUnit, error) {
	u, ok := mapDurationUnitShortName[unit]
	if ok {
		return u, nil
	}

	return DurationUnitNone, fmt.Errorf("unknown duration unit: %s", unit)
}
