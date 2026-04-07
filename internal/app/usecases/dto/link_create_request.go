package dto

import (
	"time"
)

type LinkCreateRequest struct {
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
