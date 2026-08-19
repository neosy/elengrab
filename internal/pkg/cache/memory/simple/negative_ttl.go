package memsimple

import "time"

const minNegativeTTL = 15 * time.Second

// minNegativeTTL defines the minimum lifetime of a negative cache entry,
// preventing negative results from expiring too quickly when a short
// positive TTL is configured.
type NegativeTTL func(positiveTTL time.Duration) time.Duration

var (
	// AggressiveNegativeTTL returns positiveTTL / 12, but not less than 15 seconds.
	AggressiveNegativeTTL NegativeTTL = func(ttl time.Duration) time.Duration {
		return calculateNegativeTTL(ttl, 12)
	}

	// DefaultNegativeTTL returns positiveTTL / 8, but not less than 15 seconds.
	DefaultNegativeTTL NegativeTTL = func(ttl time.Duration) time.Duration {
		return calculateNegativeTTL(ttl, 8)
	}

	// ConservativeNegativeTTL returns positiveTTL / 5, but not less than 15 seconds.
	// Use when you want negative cache to live relatively longer.
	ConservativeNegativeTTL NegativeTTL = func(ttl time.Duration) time.Duration {
		return calculateNegativeTTL(ttl, 5)
	}
)

// calculateNegativeTTL computes negative TTL as positiveTTL / divisor,
// but never less than minNegativeTTL.
func calculateNegativeTTL(positiveTTL time.Duration, divisor int) time.Duration {
	if positiveTTL <= 0 {
		return minNegativeTTL
	}

	negTTL := positiveTTL / time.Duration(divisor)
	if negTTL < minNegativeTTL {
		return minNegativeTTL
	}
	return negTTL
}
