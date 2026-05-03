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
	mediaHostStringMap = map[MediaHost]string{
		MediaHostYouTube:   "youtube",
		MediaHostFacebook:  "facebook",
		MediaHostInstagram: "instagram",
		MediaHostTwitch:    "twitch",
		MediaHostTikTok:    "tiktok",
		MediaHostVimeo:     "vimeo",
		MediaHostRutube:    "rutube",
	}

	parseMediaHostMap = map[string]MediaHost{
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
	return mediaHostStringMap[v]
}

// Ptr returns the pointer.
func (v MediaHost) Ptr() *MediaHost {
	return &v
}

// Exists returns true if the MediaHost is valid.
func (v MediaHost) Exists() bool {
	_, exists := mediaHostStringMap[v]
	return exists
}

// ParseMediaHost converting string to MediaHost
func ParseMediaHost(s string) (MediaHost, error) {
	mediaHost, exists := parseMediaHostMap[strings.ToLower(s)]
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
