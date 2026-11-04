package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type VideoFormat string

const (
	VideoFormatNone VideoFormat = ""
	VideoFormatAuto VideoFormat = "auto"
	VideoFormatMP4  VideoFormat = "mp4"
	VideoFormatWebM VideoFormat = "webm"
)

var (
	videoFormats = map[VideoFormat]struct{}{
		VideoFormatAuto: {},
		VideoFormatMP4:  {},
		VideoFormatWebM: {},
	}
	videoFormatToFileFormatMap = map[VideoFormat]FileFormat{
		VideoFormatAuto: FileFormatAuto,
		VideoFormatMP4:  FileFormatMP4,
		VideoFormatWebM: FileFormatWebM,
	}
)

// String returns the value as a string.
func (v VideoFormat) String() string {
	return string(v)
}

// FileFormat returns the corresponding FileFormat value for the VideoFormat.
func (v VideoFormat) FileFormat() FileFormat {
	return videoFormatToFileFormatMap[v]
}

// Ptr returns the pointer.
func (v VideoFormat) Ptr() *VideoFormat {
	return &v
}

// Exists checks if the VideoFormat value exists in the predefined set of video formats.
func (v VideoFormat) Exists() bool {
	_, exists := videoFormats[v]
	return exists
}

// ParseVideoFormat converting string to VideoFormat
func ParseVideoFormat(s string) (VideoFormat, error) {
	videoFormat := VideoFormat(strings.ToLower(s))
	if !videoFormat.Exists() {
		return "", errors.New("invalid value for VideoFormat")
	}
	return videoFormat, nil
}

// MustParseVideoFormat converts a string to VideoFormat and returns VideoFormatNone if the conversion fails.
func MustParseVideoFormat(s string) VideoFormat {
	videoFormat, err := ParseVideoFormat(s)
	if err != nil {
		return VideoFormatNone
	}

	return videoFormat
}

// ValidateVideoFormat checks if the field value is a valid VideoFormat enum.
func ValidateVideoFormat(fl validator.FieldLevel) bool {
	return VideoFormat(fl.Field().String()).Exists()
}
