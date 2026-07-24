package humanize

import (
	"strconv"
)

type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// CompactNumber formats a number using compact SI-style suffixes.
//
// The following suffixes are used:
//   - K for thousands
//   - M for millions
//   - B for billions
//
// Examples:
//   - 999       -> "999"
//   - 1_000     -> "1K"
//   - 15_400    -> "15K"
//   - 1_250_000 -> "1M"
func CompactNumber[T unsigned](num T) string {
	n := uint64(num)

	switch {
	case n >= 1_000_000_000:
		return strconv.FormatUint(uint64(n/1_000_000_000), 10) + "B"
	case n >= 1_000_000:
		return strconv.FormatUint(uint64(n/1_000_000), 10) + "M"
	case n >= 1_000:
		return strconv.FormatUint(uint64(n/1_000), 10) + "K"
	default:
		return strconv.FormatUint(uint64(n), 10)
	}
}
