package authdto

import "time"

type Token struct {
	Token        string
	ExpiresAt    time.Time
	NeedsRefresh bool
}
