package dlink

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	shcode "github.com/neosy/elengrab/internal/app/utils/short_code"
)

type Link struct {
	// Unique identifier of the link
	LinkID uuid.UUID

	// Original (long) URL to be shortened
	OriginalURL string

	// Generated short code used in the shortened URL, e.g., "abc123
	ShortCode string

	// Full short URL, including domain, e.g., "https://s.nhub.ru/abc123"
	ShortURL string

	// Indicates if the full short URL should be used for exact match
	IsMatchShortURL bool

	// Maximum number of allowed clicks; nil means unlimited
	MaxClicks *uint16

	// Array of user IDs allowed to access the link; nil means no restrictions
	AllowedUserIDs []string

	// Array of IP addresses allowed to access the link; nil means no restrictions
	AllowedIPs []string

	// Expiration date and time for the link; nil means no expiration
	ExpiresAt *time.Time

	// Timestamp when the link was created
	CreatedAt time.Time

	// Timestamp when the link was last updated
	UpdatedAt time.Time

	// Timestamp when the link was soft-deleted
	DeletedAt *time.Time
}

const (
	// Min length in characters in the short link code
	minLinkShortCodeLength uint8 = 4
	// Length in characters in the short link code (ShortCode)
	linkShortCodeLengthDefault uint8 = 8
)

// GenerateShortCode generates a short code for a link based on the link ID, URL, and current timestamp.
// deduplicate - if true, ensures that the same original URL always maps to a single short code;
// if a short link for the given URL already exists, it will be returned instead of creating a new one.
// if false, a new unique short code is generated on each request, even for identical URLs.
func GenerateShortCode(linkID uuid.UUID, url string, len uint8, deduplicate bool) string {
	// Set default value for short code length
	shortCodeLength := linkShortCodeLengthDefault

	// If a length (len) is provided and greater than zero — use it
	if len > 0 {
		shortCodeLength = len
	}

	// Minimum allowed length is minLinkShortCodeLength; if less is provided — enforce the minimum
	if shortCodeLength < minLinkShortCodeLength {
		shortCodeLength = minLinkShortCodeLength
	}

	// Combine UUID, URL, and timestamp into a single string
	var combined string
	if deduplicate {
		combined = url
	} else {
		combined = fmt.Sprintf("%s:%s:%d", linkID.String(), url, time.Now().UTC().UnixNano())
	}

	// Generate short code using the combined string
	return shcode.GenerateShortCode(combined, shortCodeLength)
}

// GetShortCodeFromURL extracts and returns the last segment of the ShortCode URL path.
// For example, if ShortCode is "https://example.com/abc123", it returns "abc123".
func GetShortCodeFromURL(shortURL string) string {
	parts := strings.Split(strings.TrimRight(shortURL, "/"), "/")

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
