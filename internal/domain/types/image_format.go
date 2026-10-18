package dtypes

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ImageFormat uint8

const (
	ImageFormatUnknown ImageFormat = iota
	ImageFormatJPEG
	ImageFormatPNG
	ImageFormatWebP
	ImageFormatIcon
	ImageFormatSVG
)

var (
	imageFormatStringMap = map[ImageFormat]string{
		ImageFormatJPEG: "jpg",
		ImageFormatPNG:  "png",
		ImageFormatWebP: "webp",
		ImageFormatIcon: "ico",
		ImageFormatSVG:  "svg",
	}

	parseImageFormatMap = map[string]ImageFormat{
		"jpg":  ImageFormatJPEG,
		"jfif": ImageFormatJPEG,
		"png":  ImageFormatPNG,
		"webp": ImageFormatWebP,
		"ico":  ImageFormatIcon,
		"svg":  ImageFormatSVG,
	}
)

// String returns the value as a string.
func (v ImageFormat) String() string {
	return imageFormatStringMap[v]
}

// Ptr returns the pointer.
func (v ImageFormat) Ptr() *ImageFormat {
	return &v
}

// Exists returns true if the ImageFormat is valid.
func (v ImageFormat) Exists() bool {
	_, exists := imageFormatStringMap[v]
	return exists
}

// ParseImageFormat converting string to ImageFormat
func ParseImageFormat(s string) (ImageFormat, error) {
	imageFormat, exists := parseImageFormatMap[strings.ToLower(s)]
	if !exists {
		return ImageFormatUnknown, errors.New("invalid value for ImageFormat")
	}
	return imageFormat, nil
}

// MustParseImageFormat converts string to ImageFormat, returning ImageFormatUnknown if invalid.
func MustParseImageFormat(s string) ImageFormat {
	imageFormat, err := ParseImageFormat(s)
	if err != nil {
		return ImageFormatUnknown
	}
	return imageFormat
}

// ImageFormatFromFileName converting ext from fileName to ImageFormat
func ImageFormatFromFileName(fileName string) ImageFormat {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return ImageFormatUnknown
	}
	return MustParseImageFormat(strings.TrimPrefix(ext, "."))
}

// ValidateImageFormat checks if the field value is a valid ImageFormat enum.
func ValidateImageFormat(fl validator.FieldLevel) bool {
	_, err := ParseImageFormat(fl.Field().String())
	return err == nil
}
