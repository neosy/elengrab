package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MediaHost uint8

const (
	MediaHostNone MediaHost = iota
	MediaHostYouTube
	MediaHostFacebook
	MediaHostInstagram
	MediaHostTwitch
	MediaHostTikTok
	MediaHostVimeo
	MediaHostRutube
)

var (
	mapMediaHostString = map[MediaHost]string{
		MediaHostYouTube:   "youtube",
		MediaHostFacebook:  "facebook",
		MediaHostInstagram: "instagram",
		MediaHostTwitch:    "twitch",
		MediaHostTikTok:    "tiktok",
		MediaHostVimeo:     "vimeo",
		MediaHostRutube:    "rutube",
	}

	mapMediaHostTitle = map[MediaHost]string{
		MediaHostYouTube:   "YouTube",
		MediaHostFacebook:  "Facebook",
		MediaHostInstagram: "Instagram",
		MediaHostTwitch:    "Twitch",
		MediaHostTikTok:    "TikTok",
		MediaHostVimeo:     "Vimeo",
		MediaHostRutube:    "RuTube",
	}

	mapParseMediaHost = map[string]MediaHost{
		"youtube":   MediaHostYouTube,
		"facebook":  MediaHostFacebook,
		"instagram": MediaHostInstagram,
		"twitch":    MediaHostTwitch,
		"tiktok":    MediaHostTikTok,
		"vimeo":     MediaHostVimeo,
		"rutube":    MediaHostRutube,
	}
)

// String returns the value as a string.
func (v MediaHost) String() string {
	return mapMediaHostString[v]
}

// Title returns the title.
func (v MediaHost) Title() string {
	return mapMediaHostTitle[v]
}

// Ptr returns the pointer.
func (v MediaHost) Ptr() *MediaHost {
	return &v
}

// Exists returns true if the MediaHost is valid.
func (v MediaHost) Exists() bool {
	_, exists := mapMediaHostString[v]
	return exists
}

// ParseMediaHost converting string to MediaHost
func ParseMediaHost(s string) (MediaHost, error) {
	mediaHost, exists := mapParseMediaHost[strings.ToLower(s)]
	if !exists {
		return MediaHostNone, errors.New("invalid value for MediaHost")
	}
	return mediaHost, nil
}

// ValidateMediaHost checks if the field value is a valid MediaHost enum.
func ValidateMediaHost(fl validator.FieldLevel) bool {
	_, err := ParseMediaHost(fl.Field().String())
	return err == nil
}

func (v MediaHost) ThumbnailSourceType() ThumbnailSourceType {
	return MapMediaHostToThumbnailSourceType(v)
}
