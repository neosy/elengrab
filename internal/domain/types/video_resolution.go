package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type VideoResolution string

const (
	VideoResolutionNone  VideoResolution = "none"
	VideoResolutionBest  VideoResolution = "best"
	VideoResolution4k    VideoResolution = "4k"
	VideoResolution2k    VideoResolution = "2k"
	VideoResolution1080p VideoResolution = "1080p"
	VideoResolution720p  VideoResolution = "720p"
	VideoResolution480p  VideoResolution = "480p"
	VideoResolution360p  VideoResolution = "360p"
)

var (
	videoResolutionMap = map[VideoResolution]struct{}{
		VideoResolutionNone:  {},
		VideoResolutionBest:  {},
		VideoResolution4k:    {},
		VideoResolution2k:    {},
		VideoResolution1080p: {},
		VideoResolution720p:  {},
		VideoResolution480p:  {},
		VideoResolution360p:  {},
	}

	videoResolutionHeightMap = map[VideoResolution]uint16{
		VideoResolutionNone:  0,
		VideoResolutionBest:  0,
		VideoResolution4k:    2160,
		VideoResolution2k:    1440,
		VideoResolution1080p: 1080,
		VideoResolution720p:  720,
		VideoResolution480p:  480,
		VideoResolution360p:  360,
	}

	videoResolutionWidthMap = map[VideoResolution]uint16{
		VideoResolutionNone:  0,
		VideoResolutionBest:  0,
		VideoResolution4k:    3840,
		VideoResolution2k:    2560,
		VideoResolution1080p: 1920,
		VideoResolution720p:  1280,
		VideoResolution480p:  854,
		VideoResolution360p:  640,
	}
)

// String returns the value as a string.
func (v VideoResolution) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v VideoResolution) Ptr() *VideoResolution {
	return &v
}

// Exists returns true if the VideoResolution is valid.
func (v VideoResolution) Exists() bool {
	_, exists := videoResolutionMap[v]
	return exists
}

// ParseVideoResolution converting string to VideoResolution
func ParseVideoResolution(s string) (VideoResolution, error) {
	videoResolution := VideoResolution(strings.ToLower(s))

	if _, exists := videoResolutionMap[videoResolution]; !exists {
		return "", errors.New("invalid value for VideoResolution")
	}

	return videoResolution, nil
}

// ValidateVideoResolution checks if the field value is a valid VideoResolution enum.
func ValidateVideoResolution(fl validator.FieldLevel) bool {
	return VideoResolution(fl.Field().String()).Exists()
}

func (v VideoResolution) Height() uint16 {
	return videoResolutionHeightMap[v]
}

func (v VideoResolution) Width() uint16 {
	return videoResolutionWidthMap[v]
}
