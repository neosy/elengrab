package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MediaPlatform uint8

const (
	MediaPlatformNone MediaPlatform = iota
	MediaPlatformYouTube
	MediaPlatformFacebook
	MediaPlatformInstagram
	MediaPlatformTwitch
	MediaPlatformTikTok
	MediaPlatformVimeo
	MediaPlatformRutube
)

var (
	mapMediaPlatformString = map[MediaPlatform]string{
		MediaPlatformYouTube:   "youtube",
		MediaPlatformFacebook:  "facebook",
		MediaPlatformInstagram: "instagram",
		MediaPlatformTwitch:    "twitch",
		MediaPlatformTikTok:    "tiktok",
		MediaPlatformVimeo:     "vimeo",
		MediaPlatformRutube:    "rutube",
	}

	mapMediaPlatformTitle = map[MediaPlatform]string{
		MediaPlatformYouTube:   "YouTube",
		MediaPlatformFacebook:  "Facebook",
		MediaPlatformInstagram: "Instagram",
		MediaPlatformTwitch:    "Twitch",
		MediaPlatformTikTok:    "TikTok",
		MediaPlatformVimeo:     "Vimeo",
		MediaPlatformRutube:    "RuTube",
	}

	mapParseMediaPlatform = map[string]MediaPlatform{
		"youtube":   MediaPlatformYouTube,
		"facebook":  MediaPlatformFacebook,
		"instagram": MediaPlatformInstagram,
		"twitch":    MediaPlatformTwitch,
		"tiktok":    MediaPlatformTikTok,
		"vimeo":     MediaPlatformVimeo,
		"rutube":    MediaPlatformRutube,
	}
)

// String returns the value as a string.
func (v MediaPlatform) String() string {
	return mapMediaPlatformString[v]
}

// Title returns the title.
func (v MediaPlatform) Title() string {
	return mapMediaPlatformTitle[v]
}

// Ptr returns the pointer.
func (v MediaPlatform) Ptr() *MediaPlatform {
	return &v
}

// Exists returns true if the MediaPlatform is valid.
func (v MediaPlatform) Exists() bool {
	_, exists := mapMediaPlatformString[v]
	return exists
}

// ParseMediaPlatform converting string to MediaPlatform
func ParseMediaPlatform(s string) (MediaPlatform, error) {
	mediaPlatform, exists := mapParseMediaPlatform[strings.ToLower(s)]
	if !exists {
		return MediaPlatformNone, errors.New("invalid value for MediaPlatform")
	}
	return mediaPlatform, nil
}

// ValidateMediaPlatform checks if the field value is a valid MediaPlatform enum.
func ValidateMediaPlatform(fl validator.FieldLevel) bool {
	_, err := ParseMediaPlatform(fl.Field().String())
	return err == nil
}

func (v MediaPlatform) ThumbnailSourceType() ThumbnailSourceType {
	return MapMediaPlatformToThumbnailSourceType(v)
}
