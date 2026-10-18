package pservices

import "time"

type ShortLinkClickRequest struct {
	ShortURL string

	ClickedBy *string

	IPAddress string
	UserAgent *string
	Referrer  *string
}

type ShortLinkCreateRequest struct {
	OriginalURL     string
	BaseURL         *string
	ShortCodeLength *uint8
	IsMatchShortURL bool
	MaxClicks       *uint16
	AllowedUserIDs  []string
	AllowedIPs      []string
	ExpiresAt       *time.Time
	Deterministic   *bool
}
