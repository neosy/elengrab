package uformat

import (
	"fmt"
	"strings"
	"time"
)

type (
	// DurationFormatOptions defines configurable options for formatting a time.Duration
	// using DurationFormat. It allows customizing precision, unit display, and spacing.
	DurationFormatOptions struct {
		// ForceUnit, if set, overrides automatic unit selection and forces
		// the output to use the specified DurationUnit (e.g., Millisecond, Second).
		ForceUnit *DurationUnit

		// Precision defines the number of decimal places to display.
		// For example, Precision=2 → "6.10ms".
		Precision uint8

		// ShowUnit indicates whether to append the unit suffix (ns, µs, ms, s, etc.)
		// to the formatted output. If false, only the numeric value is shown.
		ShowUnit bool

		// SpaceBeforeUnit controls whether to insert a space between the number
		// and the unit. Example: true → "6.10 ms", false → "6.10ms".
		SpaceBeforeUnit bool

		// TrailingZero indicates whether to keep trailing zeros in the fractional part.
		// Example: true → "6.10ms", false → "6.1m".
		TrailingZero bool
	}

	// DurationFormatOption is a functional option used to modify
	// DurationFormatOptions when calling DurationFormat. This allows
	// flexible configuration without manually constructing the struct.
	DurationFormatOption func(*DurationFormatOptions)
)

// DefaultDurationFormatOptions creates a new DurationFormatOptions with default values.
func DefaultDurationFormatOptions() DurationFormatOptions {
	return DurationFormatOptions{
		ForceUnit:       nil,
		Precision:       2,
		ShowUnit:        true,
		SpaceBeforeUnit: false,
		TrailingZero:    true,
	}
}

// DurationFormatOptionForceUnit sets the force unit for the formatted duration.
func DurationFormatOptionForceUnit(unit DurationUnit) DurationFormatOption {
	return func(o *DurationFormatOptions) {
		o.ForceUnit = &unit
	}
}

// DurationFormatOptionPrecision sets the precision for the formatted duration.
func DurationFormatOptionPrecision(precision uint8) DurationFormatOption {
	return func(o *DurationFormatOptions) {
		o.Precision = precision
	}
}

// DurationFormatOptionShowUnit sets whether to show the unit in the formatted duration.
func DurationFormatOptionShowUnit(showUnit bool) DurationFormatOption {
	return func(o *DurationFormatOptions) {
		o.ShowUnit = showUnit
	}
}

// DurationFormatOptionSpaceBeforeUnit sets whether to add a space before the unit in the formatted duration.
func DurationFormatOptionSpaceBeforeUnit(spaceBeforeUnit bool) DurationFormatOption {
	return func(o *DurationFormatOptions) {
		o.SpaceBeforeUnit = spaceBeforeUnit
	}
}

// DurationFormatOptionTrailingZero sets whether to keep trailing zeros.
func DurationFormatOptionTrailingZero(trailingZero bool) DurationFormatOption {
	return func(o *DurationFormatOptions) {
		o.TrailingZero = trailingZero
	}
}

// DurationFormat formats a time.Duration into a human-readable string
// according to the specified DurationFormatOptions.
//
// The formatting behavior can be customized using functional options:
//   - Precision: number of decimal places to display.
//   - TrailingZero: whether to keep trailing zeros after the decimal point.
//   - ShowUnit: whether to append the time unit (ns, µs, ms, s, etc.).
//   - SpaceBeforeUnit: whether to insert a space between the numeric value and the unit.
//   - ForceUnit: if set, overrides automatic unit selection and forces the output
//     to use the specified DurationUnit (e.g., Millisecond, Second).
//
// Example usage:
//
//	d := 6104888 * time.Nanosecond
//
//	s := uformat.DurationFormat(d,
//	    uformat.WithPrecision(2),
//	    uformat.WithTrailingZero(false),
//	)
//	fmt.Println(s) // → "6.1ms"
//
// The function returns a string representing the duration in the chosen unit
// and formatting style.
func DurationFormat(d time.Duration, opts ...DurationFormatOption) string {
	opt := DefaultDurationFormatOptions()

	for _, o := range opts {
		o(&opt)
	}

	suffix := func(unit string) string {
		var suffix string

		if opt.ShowUnit && unit != "" {
			if opt.SpaceBeforeUnit {
				suffix += " "
			}
			suffix += unit
		}

		return suffix
	}

	var unit DurationUnit
	if opt.ForceUnit != nil {
		unit = *opt.ForceUnit
	} else {
		unit = DurationUnitByValue(d)
	}

	numStr := fmt.Sprintf("%.*f", opt.Precision, unit.ToFloat64(d))

	if !opt.TrailingZero {
		numStr = strings.TrimRight(numStr, "0")
		numStr = strings.TrimRight(numStr, ".")
	}

	return fmt.Sprintf("%s%s", numStr, suffix(unit.String()))
}
