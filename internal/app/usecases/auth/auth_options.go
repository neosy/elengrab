package auth

import "time"

const (
	// Maximum lifetime of a session (absolute expiration)
	defaultSessionTTL = 30 * 24 * time.Hour

	// Refresh interval: how often we extend the session expiration on activity
	// This means that if the user is active, their session will be extended every 5 days,
	// up to the maximum of 30 days.
	defaultSessionRefreshInterval = 5 * 24 * time.Hour
)

type (
	AuthOptions struct {
		// Maximum lifetime of a session (absolute expiration)
		SessionTTL time.Duration

		// Refresh interval: how often we extend the session expiration on activity
		SessionRefreshInterval time.Duration
	}
	AuthOption func(*AuthOptions)
)

var (
	defaultAuthOptions = AuthOptions{
		SessionTTL:             defaultSessionTTL,
		SessionRefreshInterval: defaultSessionRefreshInterval,
	}
)

func NewAuthOptions(opts ...AuthOption) AuthOptions {
	options := defaultAuthOptions
	options.init(opts...)
	return options
}

func (o *AuthOptions) init(opts ...AuthOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithSessionTTL sets the maximum lifetime of a session (absolute expiration).
// This means that regardless of user activity, the session will expire after this duration.
func WithSessionTTL(value time.Duration) AuthOption {
	return func(options *AuthOptions) {
		if value > 0 {
			options.SessionTTL = value
		} else {
			options.SessionTTL = defaultSessionTTL
		}
	}
}

// WithSessionRefreshInterval sets the refresh interval for extending session expiration on activity.
// This means that if the user is active, their session will be extended every specified duration,
// up to the maximum session TTL.
func WithSessionRefreshInterval(value time.Duration) AuthOption {
	return func(options *AuthOptions) {
		if value > 0 {
			options.SessionRefreshInterval = value
		} else {
			options.SessionRefreshInterval = defaultSessionRefreshInterval
		}
	}
}
