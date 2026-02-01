package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AudioCodec string

const (
	AudioCodecNone AudioCodec = "none"
	AudioCodecMP3  AudioCodec = "mp3"
	AudioCodecAAC  AudioCodec = "aac"
	AudioCodecFLAC AudioCodec = "flac"
	AudioCodecOPUS AudioCodec = "opus"
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

	audioCodecToAudioFormatMap = map[AudioCodec]AudioFormat{
		AudioCodecMP3:  AudioFormatMP3,
		AudioCodecAAC:  AudioFormatM4A,
		AudioCodecFLAC: AudioFormatFLAC,
		AudioCodecOPUS: AudioFormatOPUS,
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

// AudioFormat returns the corresponding AudioFormat for the AudioCodec.
func (v AudioCodec) AudioFormat() AudioFormat {
	format, exists := audioCodecToAudioFormatMap[v]
	if !exists {
		return AudioFormatNone
	}
	return format
}

// ParseAudioCodec converting string to AudioCodec
func ParseAudioCodec(s string) (AudioCodec, error) {
	audioCodec := AudioCodec(strings.ToLower(s))

	if _, exists := audioCodecMap[audioCodec]; !exists {
		return "", errors.New("invalid value for AudioCodec")
	}

	return audioCodec, nil
}

// MustParseAudioCodec
func MustParseAudioCodec(s string) AudioCodec {
	codec, err := ParseAudioCodec(s)
	if err != nil {
		return AudioCodecNone
	}
	return codec
}

// ValidateAudioCodec checks if the field value is a valid AudioCodec enum.
func ValidateAudioCodec(fl validator.FieldLevel) bool {
	return AudioCodec(fl.Field().String()).Exists()
}
