package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MediaWatchEventType string

const (
	MediaWatchEventTypeNone      MediaWatchEventType = ""
	MediaWatchEventTypePause     MediaWatchEventType = "pause"
	MediaWatchEventTypeEnded     MediaWatchEventType = "ended"
	MediaWatchEventTypeSeek      MediaWatchEventType = "seek"
	MediaWatchEventTypeHeartbeat MediaWatchEventType = "heartbeat"
)

var (
	mediaWatchEventTypeStringMap = map[MediaWatchEventType]string{
		MediaWatchEventTypePause:     "pause",
		MediaWatchEventTypeEnded:     "ended",
		MediaWatchEventTypeSeek:      "seek",
		MediaWatchEventTypeHeartbeat: "heartbeat",
	}

	parseMediaWatchEventTypeMap = map[string]MediaWatchEventType{
		"pause":     MediaWatchEventTypePause,
		"ended":     MediaWatchEventTypeEnded,
		"seek":      MediaWatchEventTypeSeek,
		"heartbeat": MediaWatchEventTypeHeartbeat,
	}
)

// String returns the value as a string.
func (v MediaWatchEventType) String() string {
	return mediaWatchEventTypeStringMap[v]
}

// Ptr returns the pointer.
func (v MediaWatchEventType) Ptr() *MediaWatchEventType {
	return &v
}

// Exists returns true if the MediaWatchEventType is valid.
func (v MediaWatchEventType) Exists() bool {
	_, exists := mediaWatchEventTypeStringMap[v]
	return exists
}

// ParseMediaWatchEventType converting string to MediaWatchEventType
func ParseMediaWatchEventType(s string) (MediaWatchEventType, error) {
	mediaWatchEventType, exists := parseMediaWatchEventTypeMap[strings.ToLower(s)]
	if !exists {
		return MediaWatchEventTypeNone, errors.New("invalid value for MediaWatchEventType")
	}
	return mediaWatchEventType, nil
}

// ValidateMediaWatchEventType checks if the field value is a valid MediaWatchEventType enum.
func ValidateMediaWatchEventType(fl validator.FieldLevel) bool {
	_, err := ParseMediaWatchEventType(fl.Field().String())
	return err == nil
}
