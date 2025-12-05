package nworkers

import "time"

// Interval represents a duration with a default value.
type Interval struct {
	vDef time.Duration
	v    time.Duration
}

// NewInterval creates a new Interval with a value and a default.
// If value is zero, default will be used.
func NewInterval(valueDefault time.Duration, value time.Duration) Interval {
	return Interval{
		vDef: valueDefault,
		v:    value,
	}
}

// Value returns the effective duration: the actual value if set, otherwise the default.
func (i *Interval) Duration() time.Duration {
	if i.v.Seconds() >= 1 {
		return i.v
	}
	return i.vDef
}

// DurationPtr returns a pointer to the effective duration.
func (i *Interval) DurationPtr() *time.Duration {
	v := i.Duration()
	return &v
}
