package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type VideoFormat string

const (
	VideoFormatNone VideoFormat = "none"
	VideoFormatAuto VideoFormat = "auto"
	VideoFormatMP4  VideoFormat = "mp4"
	VideoFormatWebM VideoFormat = "webm"
)

var (
	videoFormatMap = map[VideoFormat]struct{}{
		VideoFormatNone: {},
		VideoFormatAuto: {},
		VideoFormatMP4:  {},
		VideoFormatWebM: {},
	}
	parseVideoFormatMap = map[string]VideoFormat{
		"none": VideoFormatNone,
		"auto": VideoFormatAuto,
		"mp4":  VideoFormatMP4,
		"webm": VideoFormatWebM,
	}
	videoFormatToFileFormatMap = map[VideoFormat]FileFormat{
		VideoFormatNone: FileFormatNone,
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

// Exists returns true if the VideoFormat is valid.
func (v VideoFormat) Exists() bool {
	_, exists := videoFormatMap[v]
	return exists
}

// ParseVideoFormat converting string to VideoFormat
func ParseVideoFormat(s string) (VideoFormat, error) {
	videoFormat, exists := parseVideoFormatMap[strings.ToLower(s)]
	if !exists {
		return "", errors.New("invalid value for VideoFormat")
	}
	return videoFormat, nil
}

// ValidateVideoFormat checks if the field value is a valid VideoFormat enum.
func ValidateVideoFormat(fl validator.FieldLevel) bool {
	return VideoFormat(fl.Field().String()).Exists()
}
