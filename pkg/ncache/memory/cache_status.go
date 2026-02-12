package nmemory

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CacheStatus represents the status of a cache entry.
type CacheStatus uint8

const (
	// CacheStatusMiss indicates that the cache has no entry for the channel,
	// so a DB lookup is required.
	CacheStatusMiss CacheStatus = iota

	// CacheStatusHit indicates that the channel was found in the cache.
	CacheStatusHit

	// CacheStatusNegativeHit indicates that the cache contains a negative entry,
	// meaning the channel does not exist (negative cache).
	CacheStatusNegativeHit
)

var (
	// cacheStatusMap implementation of a set for CacheStatus
	cacheStatusMap = map[CacheStatus]string{
		CacheStatusMiss:        "miss",
		CacheStatusHit:         "hit",
		CacheStatusNegativeHit: "negative",
	}

	parseCacheStatusMap = map[string]CacheStatus{
		"miss":     CacheStatusMiss,
		"hit":      CacheStatusHit,
		"negative": CacheStatusNegativeHit,
	}
)

// String returns the value as a string.
func (v CacheStatus) String() string {
	return cacheStatusMap[v]
}

// Exists returns true if the CacheStatus is valid.
func (v CacheStatus) Exists() bool {
	_, exists := cacheStatusMap[v]
	return exists
}

// ParseCacheStatus converting string to CacheStatus
func ParseCacheStatus(s string) (CacheStatus, error) {
	status, exists := parseCacheStatusMap[strings.ToLower(s)]
	if !exists {
		return CacheStatusMiss, errors.New("invalid value for CacheStatus")
	}
	return status, nil
}

// MustParseCacheStatus converting string to CacheStatus, ignoring any errors.
func MustParseCacheStatus(s string) CacheStatus {
	status, _ := ParseCacheStatus(s)
	return status
}

// ValidateCacheStatus checks if the field value is a valid CacheStatus enum.
func ValidateCacheStatus(fl validator.FieldLevel) bool {
	_, exist := parseCacheStatusMap[fl.Field().String()]
	return exist
}
