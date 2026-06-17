package dtypes

import (
	"errors"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/neosy/elengrab/internal/pkg/stringx"
)

// Visibility access level for media
type MediaVisibility uint8

const (
	// Media is private and accessible only with permissions
	MediaVisibilityPrivate MediaVisibility = iota

	// Media is accessible to any authenticated (logged-in) user in the system,
	// but not to anonymous visitors.
	MediaVisibilityAuthenticated

	// Media is publicly accessible
	MediaVisibilityPublic

	// Default visibility for new media
	MediaVisibilityDefault = MediaVisibilityPrivate
)

var (
	// nameByMediaVisibility maps an MediaVisibility to its string representation
	nameByMediaVisibility = map[MediaVisibility]string{
		MediaVisibilityPrivate:       "private",
		MediaVisibilityAuthenticated: "authenticated",
		MediaVisibilityPublic:        "public",
	}

	mediaVisibilityByName = map[string]MediaVisibility{
		"private":       MediaVisibilityPrivate,
		"authenticated": MediaVisibilityAuthenticated,
		"public":        MediaVisibilityPublic,
	}

	mediaVisibilityList = []MediaVisibility{
		MediaVisibilityPrivate,
		MediaVisibilityAuthenticated,
		MediaVisibilityPublic,
	}

	mediaVisibilityByAppMode = map[AppMode]MediaVisibility{
		AppModePublic: MediaVisibilityPublic,
	}
)

// String returns the value as a string.
func (v MediaVisibility) String() string {
	return nameByMediaVisibility[v]
}

// Label returns the value as a label.
func (v MediaVisibility) Label() string {
	return stringx.Capitalize(v.String())
}

// Exists returns true if the MediaVisibility is valid.
func (v MediaVisibility) Exists() bool {
	_, exists := nameByMediaVisibility[v]
	return exists
}

// ParseMediaVisibility converting string to MediaVisibility
func ParseMediaVisibility(s string) (MediaVisibility, error) {
	mode, exists := mediaVisibilityByName[strings.ToLower(s)]
	if !exists {
		return MediaVisibilityDefault, errors.New("invalid value for MediaVisibility")
	}
	return mode, nil
}

// MustParseMediaVisibility converting string to MediaVisibility, ignoring any errors.
func MustParseMediaVisibility(s string) MediaVisibility {
	mode, err := ParseMediaVisibility(s)
	if err != nil {
		return MediaVisibilityDefault
	}
	return mode
}

// List
func MediaVisibilityList() []MediaVisibility {
	return slices.Clone(mediaVisibilityList)
}

func MediaVisibilityByAppMode(appMode AppMode) MediaVisibility {
	visibility, exists := mediaVisibilityByAppMode[appMode]
	if exists {
		return visibility
	}
	return MediaVisibilityDefault
}

// ValidateMediaVisibility checks if the field value is a valid MediaVisibility enum.
func ValidateMediaVisibility(fl validator.FieldLevel) bool {
	_, exist := mediaVisibilityByName[fl.Field().String()]
	return exist
}
