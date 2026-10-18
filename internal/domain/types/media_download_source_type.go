package dtypes

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MediaDownloadSourceType uint8

const (
	MediaDownloadSourceTypeNone MediaDownloadSourceType = iota
	MediaDownloadSourceTypeURL
	MediaDownloadSourceTypeUpload
)

var (
	mediaDownloadSourceTypes = map[MediaDownloadSourceType]string{
		MediaDownloadSourceTypeURL:    "url",
		MediaDownloadSourceTypeUpload: "upload",
	}

	mediaDownloadSourceTypesByStringValue = map[string]MediaDownloadSourceType{
		"url":    MediaDownloadSourceTypeURL,
		"upload": MediaDownloadSourceTypeUpload,
	}
)

// String returns the value as a string.
func (v MediaDownloadSourceType) String() string {
	return mediaDownloadSourceTypes[v]
}

// Validate checks that the source type is one of the supported values.
func (v MediaDownloadSourceType) Validate() error {
	if !v.Exists() {
		return fmt.Errorf("invalid media download source type: %v", v)
	}
	return nil
}

// Ptr returns the pointer.
func (v MediaDownloadSourceType) Ptr() *MediaDownloadSourceType {
	return &v
}

// Exists returns true if the MediaDownloadSourceType is valid.
func (v MediaDownloadSourceType) Exists() bool {
	_, exists := mediaDownloadSourceTypes[v]
	return exists
}

// ParseMediaDownloadSourceType converting string to MediaDownloadSourceType
func ParseMediaDownloadSourceType(s string) (MediaDownloadSourceType, error) {
	mediaDownloadSourceType, exists := mediaDownloadSourceTypesByStringValue[strings.ToLower(s)]
	if !exists {
		return MediaDownloadSourceTypeNone, fmt.Errorf("invalid value for MediaDownloadSourceType: %q", s)
	}
	return mediaDownloadSourceType, nil
}

// ValidateMediaDownloadSourceType checks if the field value is a valid MediaDownloadSourceType enum.
func ValidateMediaDownloadSourceType(fl validator.FieldLevel) bool {
	_, err := ParseMediaDownloadSourceType(fl.Field().String())
	return err == nil
}
