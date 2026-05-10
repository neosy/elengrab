package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ThumbnailSourceType uint8

const (
	ThumbnailSourceTypeNone ThumbnailSourceType = iota
	ThumbnailSourceTypeYouTube
	ThumbnailSourceTypeVimeo
	ThumbnailSourceTypeInstagram
	ThumbnailSourceTypeTwitch
	ThumbnailSourceTypeTikTok
	ThumbnailSourceTypeRutube
	ThumbnailSourceTypeExternal
	ThumbnailSourceTypeVideoFrame
	ThumbnailSourceTypeGenerated
	ThumbnailSourceTypeUpload
)

var (
	thumbnailSourceTypeStringMap = map[ThumbnailSourceType]string{
		ThumbnailSourceTypeYouTube:    "youtube",
		ThumbnailSourceTypeVimeo:      "vimeo",
		ThumbnailSourceTypeInstagram:  "instagram",
		ThumbnailSourceTypeTwitch:     "twitch",
		ThumbnailSourceTypeTikTok:     "tiktok",
		ThumbnailSourceTypeRutube:     "rutube",
		ThumbnailSourceTypeExternal:   "external",
		ThumbnailSourceTypeVideoFrame: "video_frame",
		ThumbnailSourceTypeGenerated:  "generated",
		ThumbnailSourceTypeUpload:     "upload",
	}

	parseThumbnailSourceTypeMap = map[string]ThumbnailSourceType{
		"youtube":     ThumbnailSourceTypeYouTube,
		"vimeo":       ThumbnailSourceTypeVimeo,
		"instagram":   ThumbnailSourceTypeInstagram,
		"twitch":      ThumbnailSourceTypeTwitch,
		"tiktok":      ThumbnailSourceTypeTikTok,
		"rutube":      ThumbnailSourceTypeRutube,
		"external":    ThumbnailSourceTypeExternal,
		"video_frame": ThumbnailSourceTypeVideoFrame,
		"generated":   ThumbnailSourceTypeGenerated,
		"upload":      ThumbnailSourceTypeUpload,
	}

	mapMediaHostToThumbnailSourceType = map[MediaHost]ThumbnailSourceType{
		MediaHostYouTube:   ThumbnailSourceTypeYouTube,
		MediaHostVimeo:     ThumbnailSourceTypeVimeo,
		MediaHostInstagram: ThumbnailSourceTypeInstagram,
		MediaHostTwitch:    ThumbnailSourceTypeTwitch,
		MediaHostTikTok:    ThumbnailSourceTypeTikTok,
		MediaHostRutube:    ThumbnailSourceTypeRutube,
	}
)

// String returns the value as a string.
func (v ThumbnailSourceType) String() string {
	return thumbnailSourceTypeStringMap[v]
}

// Ptr returns the pointer.
func (v ThumbnailSourceType) Ptr() *ThumbnailSourceType {
	return &v
}

// Exists returns true if the ThumbnailSourceType is valid.
func (v ThumbnailSourceType) Exists() bool {
	_, exists := thumbnailSourceTypeStringMap[v]
	return exists
}

// ParseThumbnailSourceType converting string to ThumbnailSourceType
func ParseThumbnailSourceType(s string) (ThumbnailSourceType, error) {
	thumbnailSourceType, exists := parseThumbnailSourceTypeMap[strings.ToLower(s)]
	if !exists {
		return ThumbnailSourceTypeNone, errors.New("invalid value for ThumbnailSourceType")
	}
	return thumbnailSourceType, nil
}

// ValidateThumbnailSourceType checks if the field value is a valid ThumbnailSourceType enum.
func ValidateThumbnailSourceType(fl validator.FieldLevel) bool {
	_, err := ParseThumbnailSourceType(fl.Field().String())
	return err == nil
}

func MapMediaHostToThumbnailSourceType(mediaHost MediaHost) ThumbnailSourceType {
	return mapMediaHostToThumbnailSourceType[mediaHost]
}
