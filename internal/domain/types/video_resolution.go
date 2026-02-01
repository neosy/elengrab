package dtypes

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type VideoResolution string

const (
	VideoResolutionNone  VideoResolution = "none"
	VideoResolutionMax   VideoResolution = "max"
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
		VideoResolutionMax:   {},
		VideoResolution4k:    {},
		VideoResolution2k:    {},
		VideoResolution1080p: {},
		VideoResolution720p:  {},
		VideoResolution480p:  {},
		VideoResolution360p:  {},
	}

	videoResolutionHeightMap = map[VideoResolution]uint16{
		VideoResolutionNone:  0,
		VideoResolutionMax:   0,
		VideoResolution4k:    2160,
		VideoResolution2k:    1440,
		VideoResolution1080p: 1080,
		VideoResolution720p:  720,
		VideoResolution480p:  480,
		VideoResolution360p:  360,
	}

	videoResolutionWidthMap = map[VideoResolution]uint16{
		VideoResolutionNone:  0,
		VideoResolutionMax:   0,
		VideoResolution4k:    3840,
		VideoResolution2k:    2560,
		VideoResolution1080p: 1920,
		VideoResolution720p:  1280,
		VideoResolution480p:  854,
		VideoResolution360p:  640,
	}

	heightToResolutionMap = map[uint16]VideoResolution{
		2160: VideoResolution4k,
		1440: VideoResolution2k,
		1080: VideoResolution1080p,
		720:  VideoResolution720p,
		480:  VideoResolution480p,
		360:  VideoResolution360p,
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

// ParseVideoResolutionWH
func ParseVideoResolutionWH(w, h uint16) VideoResolution {
	if w != 0 && h > w {
		h = w
	}
	resolution, exists := heightToResolutionMap[h]
	if !exists {
		resolution = VideoResolutionNone
	}
	return resolution
}

// VideoResolutionStringToWH parses a resolution string in the form "WIDTHxHEIGHT"
// and returns width and height as uint16. Returns an error if the string is not valid.
func VideoResolutionStringToWH(res string) (uint16, uint16, error) {
	re := regexp.MustCompile(`^\d+x\d+$`)
	if !re.MatchString(res) {
		return 0, 0, fmt.Errorf("invalid video resolution format: %q", res)
	}

	values := strings.SplitN(res, "x", 2)
	if len(values) != 2 {
		return 0, 0, fmt.Errorf("cannot split resolution string: %q", res)
	}

	v1, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width in resolution: %q, %w", values[0], err)
	}
	v2, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height in resolution: %q, %w", values[1], err)
	}

	return uint16(v1), uint16(v2), nil
}

// ParseVideoResolutionFromString parses a resolution string in the form "WIDTHxHEIGHT"
// and returns a VideoResolution. Returns an error if the string is invalid.
func ParseVideoResolutionFromString(res string) (VideoResolution, error) {
	w, h, err := VideoResolutionStringToWH(res)
	if err != nil {
		return "", fmt.Errorf("invalid value for VideoResolution: %w", err)
	}

	return ParseVideoResolutionWH(w, h), nil
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
