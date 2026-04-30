package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ThumbnailVariant uint8

const (
	ThumbnailVariantNone ThumbnailVariant = iota
	ThumbnailVariantOriginal
	ThumbnailVariantSmall
	ThumbnailVariantMedium
	ThumbnailVariantLarge
)

var (
	thumbnailVariantStringMap = map[ThumbnailVariant]string{
		ThumbnailVariantOriginal: "original",
		ThumbnailVariantSmall:    "small",
		ThumbnailVariantMedium:   "medium",
		ThumbnailVariantLarge:    "large",
	}

	parseThumbnailVariantMap = map[string]ThumbnailVariant{
		"original": ThumbnailVariantOriginal,
		"small":    ThumbnailVariantSmall,
		"medium":   ThumbnailVariantMedium,
		"large":    ThumbnailVariantLarge,
	}
)

// String returns the value as a string.
func (v ThumbnailVariant) String() string {
	return thumbnailVariantStringMap[v]
}

// Ptr returns the pointer.
func (v ThumbnailVariant) Ptr() *ThumbnailVariant {
	return &v
}

// Exists returns true if the ThumbnailVariant is valid.
func (v ThumbnailVariant) Exists() bool {
	_, exists := thumbnailVariantStringMap[v]
	return exists
}

// ParseThumbnailVariant converting string to ThumbnailVariant
func ParseThumbnailVariant(s string) (ThumbnailVariant, error) {
	thumbnailVariant, exists := parseThumbnailVariantMap[strings.ToLower(s)]
	if !exists {
		return ThumbnailVariantNone, errors.New("invalid value for ThumbnailVariant")
	}
	return thumbnailVariant, nil
}

// ValidateThumbnailVariant checks if the field value is a valid ThumbnailVariant enum.
func ValidateThumbnailVariant(fl validator.FieldLevel) bool {
	_, err := ParseThumbnailVariant(fl.Field().String())
	return err == nil
}
