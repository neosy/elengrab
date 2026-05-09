package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ImageSource uint

const (
	ImageSourceNone ImageSource = iota
	ImageSourceThumbnail
	ImageSourceAvatar
	ImageSourceSite

	ImageSourceDefault = ImageSourceSite
)

var (
	imageSourceMap = map[ImageSource]string{
		ImageSourceThumbnail: "thumbnail",
		ImageSourceAvatar:    "avatar",
		ImageSourceSite:      "site",
	}

	parseImageSourceMap = map[string]ImageSource{
		"thumbnail": ImageSourceThumbnail,
		"avatar":    ImageSourceAvatar,
		"site":      ImageSourceSite,
	}
)

// String returns the value as a string.
func (v ImageSource) String() string {
	return imageSourceMap[v]
}

// Ptr returns the pointer.
func (v ImageSource) Ptr() *ImageSource {
	return &v
}

// Exists returns true if the ImageSource is valid.
func (v ImageSource) Exists() bool {
	_, exists := imageSourceMap[v]
	return exists
}

// ParseImageSource converting string to ImageSource
func ParseImageSource(s string) (ImageSource, error) {
	imageSource, exists := parseImageSourceMap[strings.ToLower(s)]
	if !exists {
		return ImageSourceNone, errors.New("invalid value for ImageSource")
	}
	return imageSource, nil
}

func MustParseImageSource(s string) ImageSource {
	imageSource, _ := ParseImageSource(s)
	return imageSource
}

// ValidateImageSource checks if the field value is a valid ImageSource enum.
func ValidateImageSource(fl validator.FieldLevel) bool {
	imageSource, _ := ParseImageSource(fl.Field().String())
	return imageSource.Exists()
}
