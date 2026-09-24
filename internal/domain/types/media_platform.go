package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MediaPlatformType uint8

const (
	MediaPlatformTypeNone MediaPlatformType = iota
	MediaPlatformTypeYouTube
	MediaPlatformTypeFacebook
	MediaPlatformTypeInstagram
	MediaPlatformTypeTwitch
	MediaPlatformTypeTikTok
	MediaPlatformTypeVimeo
	MediaPlatformTypeX
	MediaPlatformTypeRutube
	MediaPlatformTypeVKVideo
)

var (
	mapMediaPlatformTypeString = map[MediaPlatformType]string{
		MediaPlatformTypeYouTube:   "youtube",
		MediaPlatformTypeFacebook:  "facebook",
		MediaPlatformTypeInstagram: "instagram",
		MediaPlatformTypeTwitch:    "twitch",
		MediaPlatformTypeTikTok:    "tiktok",
		MediaPlatformTypeVimeo:     "vimeo",
		MediaPlatformTypeX:         "x",
		MediaPlatformTypeRutube:    "rutube",
		MediaPlatformTypeVKVideo:   "vkvideo",
	}

	mapMediaPlatformTypeTitle = map[MediaPlatformType]string{
		MediaPlatformTypeYouTube:   "YouTube",
		MediaPlatformTypeFacebook:  "Facebook",
		MediaPlatformTypeInstagram: "Instagram",
		MediaPlatformTypeTwitch:    "Twitch",
		MediaPlatformTypeTikTok:    "TikTok",
		MediaPlatformTypeVimeo:     "Vimeo",
		MediaPlatformTypeX:         "X",
		MediaPlatformTypeRutube:    "RuTube",
		MediaPlatformTypeVKVideo:   "VK Видео",
	}

	mapParseMediaPlatformType = map[string]MediaPlatformType{
		"youtube":   MediaPlatformTypeYouTube,
		"facebook":  MediaPlatformTypeFacebook,
		"instagram": MediaPlatformTypeInstagram,
		"twitch":    MediaPlatformTypeTwitch,
		"tiktok":    MediaPlatformTypeTikTok,
		"vimeo":     MediaPlatformTypeVimeo,
		"x":         MediaPlatformTypeX,
		"rutube":    MediaPlatformTypeRutube,
		"vkvideo":   MediaPlatformTypeVKVideo,
	}
)

// String returns the value as a string.
func (v MediaPlatformType) String() string {
	return mapMediaPlatformTypeString[v]
}

// Title returns the title.
func (v MediaPlatformType) Title() string {
	return mapMediaPlatformTypeTitle[v]
}

// Ptr returns the pointer.
func (v MediaPlatformType) Ptr() *MediaPlatformType {
	return &v
}

// Exists returns true if the MediaPlatformType is valid.
func (v MediaPlatformType) Exists() bool {
	_, exists := mapMediaPlatformTypeString[v]
	return exists
}

// ParseMediaPlatformType converting string to MediaPlatformType
func ParseMediaPlatformType(s string) (MediaPlatformType, error) {
	mediaPlatformType, exists := mapParseMediaPlatformType[strings.ToLower(s)]
	if !exists {
		return MediaPlatformTypeNone, errors.New("invalid value for MediaPlatformType")
	}
	return mediaPlatformType, nil
}

// ValidateMediaPlatformType checks if the field value is a valid MediaPlatformType enum.
func ValidateMediaPlatformType(fl validator.FieldLevel) bool {
	_, err := ParseMediaPlatformType(fl.Field().String())
	return err == nil
}

func (v MediaPlatformType) ThumbnailSourceType() ThumbnailSourceType {
	return MapMediaPlatformTypeToThumbnailSourceType(v)
}
