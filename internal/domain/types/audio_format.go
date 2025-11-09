package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AudioFormat string

const (
	AudioFormatNone AudioFormat = "none"
	AudioFormatOrig AudioFormat = "orig"
	AudioFormatMP3  AudioFormat = "mp3"
	AudioFormatM4A  AudioFormat = "m4a"
)

var (
	audioFormatMap = map[AudioFormat]struct{}{
		AudioFormatNone: {},
		AudioFormatOrig: {},
		AudioFormatMP3:  {},
		AudioFormatM4A:  {},
	}
)

// String returns the value as a string.
func (v AudioFormat) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v AudioFormat) Ptr() *AudioFormat {
	return &v
}

// Exists returns true if the AudioFormat is valid.
func (v AudioFormat) Exists() bool {
	_, exists := audioFormatMap[v]
	return exists
}

// ParseAudioFormat converting string to AudioFormat
func ParseAudioFormat(s string) (AudioFormat, error) {
	audioFormat := AudioFormat(strings.ToLower(s))

	if _, exists := audioFormatMap[audioFormat]; !exists {
		return "", errors.New("invalid value for AudioFormat")
	}

	return audioFormat, nil
}

// ValidateAudioFormat checks if the field value is a valid AudioFormat enum.
func ValidateAudioFormat(fl validator.FieldLevel) bool {
	return AudioFormat(fl.Field().String()).Exists()
}
