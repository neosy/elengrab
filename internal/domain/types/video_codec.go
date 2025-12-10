package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type VideoCodec string

const (
	VideoCodecNone VideoCodec = "none"
	VideoCodecBest VideoCodec = "best"
	VideoCodecH264 VideoCodec = "h264"
	VideoCodecH265 VideoCodec = "h265"
	VideoCodecAV1  VideoCodec = "av1"
	VideoCodecVP9  VideoCodec = "vp9"
)

var (
	videoCodecMap = map[VideoCodec]struct{}{
		VideoCodecNone: {},
		VideoCodecBest: {},
		VideoCodecH264: {},
		VideoCodecH265: {},
		VideoCodecAV1:  {},
		VideoCodecVP9:  {},
	}

	videoCodecTitleMap = map[VideoCodec]string{
		VideoCodecBest: "Best",
		VideoCodecH264: "H.264",
		VideoCodecH265: "H.265",
		VideoCodecAV1:  "AV1",
		VideoCodecVP9:  "VP9",
	}
)

// String returns the value as a string.
func (v VideoCodec) String() string {
	return string(v)
}

// Title
func (v VideoCodec) Title() string {
	return videoCodecTitleMap[v]
}

// Ptr returns the pointer.
func (v VideoCodec) Ptr() *VideoCodec {
	return &v
}

// Exists returns true if the VideoCodec is valid.
func (v VideoCodec) Exists() bool {
	_, exists := videoCodecMap[v]
	return exists
}

// ParseVideoCodec converting string to VideoCodec
func ParseVideoCodec(s string) (VideoCodec, error) {
	videoCodec := VideoCodec(strings.ToLower(s))

	if _, exists := videoCodecMap[videoCodec]; !exists {
		return "", errors.New("invalid value for VideoCodec")
	}

	return videoCodec, nil
}

// ValidateVideoCodec checks if the field value is a valid VideoCodec enum.
func ValidateVideoCodec(fl validator.FieldLevel) bool {
	return VideoCodec(fl.Field().String()).Exists()
}
