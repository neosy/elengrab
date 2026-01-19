package auth

import "time"

const (
	// Maximum lifetime of a session (absolute expiration)
	sessionTTL = 30 * 24 * time.Hour

	// Refresh interval: how often we extend the session expiration on activity
	sessionRefreshInterval = 5 * 24 * time.Hour
)
