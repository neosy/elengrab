package dtypes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Media task status
type MediaTaskStatus uint8

const (
	MediaTaskStatusNone MediaTaskStatus = iota
	MediaTaskStatusNew
	MediaTaskStatusPending
	MediaTaskStatusWorking
	MediaTaskStatusFailed

	MediaTaskStatusDefault = MediaTaskStatusNone
)

var (
	// MediaTaskStatuses implementation of a set for MediaTaskStatus
	mediaTaskStatuses = map[MediaTaskStatus]string{
		MediaTaskStatusNew:     "new",
		MediaTaskStatusPending: "pending",
		MediaTaskStatusWorking: "working",
		MediaTaskStatusFailed:  "failed",
	}

	mediaTaskStatusesByStringValue = map[string]MediaTaskStatus{
		"new":     MediaTaskStatusNew,
		"pending": MediaTaskStatusPending,
		"working": MediaTaskStatusWorking,
		"failed":  MediaTaskStatusFailed,
	}
)

// String returns the value as a string.
func (v MediaTaskStatus) String() string {
	return mediaTaskStatuses[v]
}

// Validate checks that the status is one of the supported values.
func (v MediaTaskStatus) Validate() error {
	if !v.Exists() {
		return fmt.Errorf("invalid media task status: %v", v)
	}
	return nil
}

// Exists returns true if the MediaTaskStatus is valid.
func (v MediaTaskStatus) Exists() bool {
	_, exists := mediaTaskStatuses[v]
	return exists
}

// ParseMediaTaskStatus converting string to status
func ParseMediaTaskStatus(s string) (MediaTaskStatus, error) {
	status, exists := mediaTaskStatusesByStringValue[strings.ToLower(s)]
	if !exists {
		return MediaTaskStatusDefault, errors.New("invalid value for MediaTaskStatus")
	}
	return status, nil
}

// MustParseMediaTaskStatus converting string to MediaTaskStatus, ignoring any errors.
func MustParseMediaTaskStatus(s string) MediaTaskStatus {
	status, _ := ParseMediaTaskStatus(s)
	return status
}

// ValidateMediaTaskStatus checks if the field value is a valid MediaTaskStatus enum.
func ValidateMediaTaskStatus(fl validator.FieldLevel) bool {
	_, exist := mediaTaskStatusesByStringValue[fl.Field().String()]
	return exist
}
