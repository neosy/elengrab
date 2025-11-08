package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FormatType string

const (
	FormatTypeNone       FormatType = "none"
	FormatTypeVideoAudio FormatType = "video_audio"
	FormatTypeVideoOnly  FormatType = "video_only"
	FormatTypeAudioOnly  FormatType = "audio_only"
)

var (
	formatTypesMap = map[FormatType]struct{}{
		FormatTypeNone:       {},
		FormatTypeVideoAudio: {},
		FormatTypeVideoOnly:  {},
		FormatTypeAudioOnly:  {},
	}
)

// String returns the value as a string.
func (v FormatType) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v FormatType) Ptr() *FormatType {
	return &v
}

// Exists returns true if the FormatType is valid.
func (v FormatType) Exists() bool {
	_, exists := formatTypesMap[v]
	return exists
}

// ParseFormatType converting string to FormatType
func ParseFormatType(s string) (FormatType, error) {
	FormatType := FormatType(strings.ToLower(s))

	if _, exists := formatTypesMap[FormatType]; !exists {
		return "", errors.New("invalid value for FormatType")
	}

	return FormatType, nil
}

// ValidateFormatType checks if the field value is a valid FormatType enum.
func ValidateFormatType(fl validator.FieldLevel) bool {
	return FormatType(fl.Field().String()).Exists()
}
