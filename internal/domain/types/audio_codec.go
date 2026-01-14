package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AudioCodec string

const (
	AudioCodecNone AudioCodec = "none"
	AudioCodecMP3  AudioCodec = "MP3"
	AudioCodecAAC  AudioCodec = "AAC"
	AudioCodecFLAC AudioCodec = "FLAC"
	AudioCodecOPUS AudioCodec = "Opus"
)

var (
	audioCodecMap = map[AudioCodec]struct{}{
		AudioCodecNone: {},
		AudioCodecMP3:  {},
		AudioCodecAAC:  {},
		AudioCodecFLAC: {},
		AudioCodecOPUS: {},
	}

	audioCodecTitleMap = map[AudioCodec]string{
		AudioCodecMP3:  "MP3",
		AudioCodecAAC:  "AAC",
		AudioCodecFLAC: "FLAC",
		AudioCodecOPUS: "Opus",
	}
)

// String returns the value as a string.
func (v AudioCodec) String() string {
	return string(v)
}

// Title
func (v AudioCodec) Title() string {
	return audioCodecTitleMap[v]
}

// Ptr returns the pointer.
func (v AudioCodec) Ptr() *AudioCodec {
	return &v
}

// Exists returns true if the AudioCodec is valid.
func (v AudioCodec) Exists() bool {
	_, exists := audioCodecMap[v]
	return exists
}

// ParseAudioCodec converting string to AudioCodec
func ParseAudioCodec(s string) (AudioCodec, error) {
	audioCodec := AudioCodec(strings.ToLower(s))

	if _, exists := audioCodecMap[audioCodec]; !exists {
		return "", errors.New("invalid value for AudioCodec")
	}

	return audioCodec, nil
}

// ValidateAudioCodec checks if the field value is a valid AudioCodec enum.
func ValidateAudioCodec(fl validator.FieldLevel) bool {
	return AudioCodec(fl.Field().String()).Exists()
}
