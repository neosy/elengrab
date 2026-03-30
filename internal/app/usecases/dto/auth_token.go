package dto

import "time"

type AuthToken struct {
	Token        string
	ExpiresAt    time.Time
	NeedsRefresh bool
}
