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
	VideoResolutionMin   VideoResolution = "min"
	VideoResolutionMax   VideoResolution = "max"
	VideoResolution16k   VideoResolution = "16k"
	VideoResolution8k    VideoResolution = "8k"
	VideoResolution4k    VideoResolution = "4k"
	VideoResolution2k    VideoResolution = "2k"
	VideoResolution1080p VideoResolution = "1080p"
	VideoResolution720p  VideoResolution = "720p"
	VideoResolution480p  VideoResolution = "480p"
	VideoResolution360p  VideoResolution = "360p"
	VideoResolution240p  VideoResolution = "240p"
	VideoResolution120p  VideoResolution = "120p"
)

var (
	videoResolutionMap = map[VideoResolution]struct{}{
		VideoResolutionNone: {},

		VideoResolutionMin: {},
		VideoResolutionMax: {},

		VideoResolution16k:   {},
		VideoResolution8k:    {},
		VideoResolution4k:    {},
		VideoResolution2k:    {},
		VideoResolution1080p: {},
		VideoResolution720p:  {},
		VideoResolution480p:  {},
		VideoResolution360p:  {},
		VideoResolution240p:  {},
		VideoResolution120p:  {},
	}

	videoResolutionOrder = []VideoResolution{
		VideoResolution120p,
		VideoResolution240p,
		VideoResolution360p,
		VideoResolution480p,
		VideoResolution720p,
		VideoResolution1080p,
		VideoResolution2k,
		VideoResolution4k,
		VideoResolution8k,
		VideoResolution16k,
	}

	videoResolutionIndex = make(map[VideoResolution]int)

	videoResolutionHeightMap = map[VideoResolution]uint16{
		VideoResolutionNone: 0,
		VideoResolutionMin:  120,
		VideoResolutionMax:  8640,

		VideoResolution16k:   8640,
		VideoResolution8k:    4320,
		VideoResolution4k:    2160,
		VideoResolution2k:    1440,
		VideoResolution1080p: 1080,
		VideoResolution720p:  720,
		VideoResolution480p:  480,
		VideoResolution360p:  360,
		VideoResolution240p:  240,
		VideoResolution120p:  120,
	}

	videoResolutionWidthMap = map[VideoResolution]uint16{
		VideoResolutionNone: 0,
		VideoResolutionMin:  160,
		VideoResolutionMax:  15360,

		VideoResolution16k:   15360,
		VideoResolution8k:    7680,
		VideoResolution4k:    3840,
		VideoResolution2k:    2560,
		VideoResolution1080p: 1920,
		VideoResolution720p:  1280,
		VideoResolution480p:  854,
		VideoResolution360p:  640,
		VideoResolution240p:  426,
		VideoResolution120p:  160,
	}

	heightToResolutionMap = map[uint16]VideoResolution{
		8640: VideoResolution16k,
		4320: VideoResolution8k,
		2160: VideoResolution4k,
		1440: VideoResolution2k,
		1080: VideoResolution1080p,
		720:  VideoResolution720p,
		480:  VideoResolution480p,
		360:  VideoResolution360p,
		240:  VideoResolution240p,
		120:  VideoResolution120p,
	}
)

func init() {
	for i, v := range videoResolutionOrder {
		videoResolutionIndex[v] = i
	}
}

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

// Height returns the height in pixels for the given VideoResolution.
func (v VideoResolution) Height() uint16 {
	return videoResolutionHeightMap[v]
}

// Width returns the width in pixels for the given VideoResolution.
func (v VideoResolution) Width() uint16 {
	return videoResolutionWidthMap[v]
}

// Prev returns the next lower VideoResolution. If there is no lower resolution, it returns VideoResolutionNone.
func (v VideoResolution) Prev() VideoResolution {
	index := videoResolutionIndex[v]
	if index == 0 || index >= len(videoResolutionOrder) {
		return VideoResolutionNone
	}

	return videoResolutionOrder[index-1]
}

// Next returns the next higher VideoResolution. If there is no higher resolution, it returns VideoResolutionNone.
func (v VideoResolution) Next() VideoResolution {
	index := videoResolutionIndex[v]
	if index-1 >= len(videoResolutionOrder) {
		return VideoResolutionNone
	}

	return videoResolutionOrder[index+1]
}
