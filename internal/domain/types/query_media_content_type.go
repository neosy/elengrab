package dtypes

import (
	"errors"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type QueryMediaContentType uint8

const (
	QueryMediaContentTypeNone QueryMediaContentType = iota
	QueryMediaContentTypeAll
	QueryMediaContentTypeVideo
	QueryMediaContentTypeShorts
	QueryMediaContentTypeAudio
	QueryMediaContentTypeInProgress
)

var (
	queryMediaContentTypes = map[QueryMediaContentType]string{
		QueryMediaContentTypeAll:        "all",
		QueryMediaContentTypeVideo:      "video",
		QueryMediaContentTypeShorts:     "shorts",
		QueryMediaContentTypeAudio:      "audio",
		QueryMediaContentTypeInProgress: "in-progress",
	}

	queryMediaContentTypeByName = map[string]QueryMediaContentType{
		"all":         QueryMediaContentTypeAll,
		"video":       QueryMediaContentTypeVideo,
		"shorts":      QueryMediaContentTypeShorts,
		"audio":       QueryMediaContentTypeAudio,
		"in-progress": QueryMediaContentTypeInProgress,
	}

	queryMediaContentTypeLabels = map[QueryMediaContentType]string{
		QueryMediaContentTypeAll:        "All",
		QueryMediaContentTypeVideo:      "Video",
		QueryMediaContentTypeShorts:     "Shorts",
		QueryMediaContentTypeAudio:      "Audio",
		QueryMediaContentTypeInProgress: "In Progress",
	}

	contentTypeList = []QueryMediaContentType{
		QueryMediaContentTypeAll,
		QueryMediaContentTypeVideo,
		QueryMediaContentTypeShorts,
		QueryMediaContentTypeAudio,
		QueryMediaContentTypeInProgress,
	}
)

// String returns the value as a string.
func (v QueryMediaContentType) String() string {
	return queryMediaContentTypes[v]
}

// Label returns the display label for the content type.
func (v QueryMediaContentType) Label() string {
	return queryMediaContentTypeLabels[v]
}

// Ptr returns the pointer.
func (v QueryMediaContentType) Ptr() *QueryMediaContentType {
	return &v
}

// Exists returns true if the QueryMediaContentType is valid.
func (v QueryMediaContentType) Exists() bool {
	_, exists := queryMediaContentTypes[v]
	return exists
}

// ParseQueryMediaContentType converting string to QueryMediaContentType
func ParseQueryMediaContentType(s string) (QueryMediaContentType, error) {
	queryMediaContentType, exists := queryMediaContentTypeByName[strings.ToLower(s)]
	if !exists {
		return QueryMediaContentTypeNone, errors.New("invalid value for QueryMediaContentType")
	}

	return queryMediaContentType, nil
}

// ValidateQueryMediaContentType checks if the field value is a valid QueryMediaContentType enum.
func ValidateQueryMediaContentType(fl validator.FieldLevel) bool {
	_, err := ParseQueryMediaContentType(fl.Field().String())
	return err == nil
}

// ContentTypeList returns all available media content types.
func ContentTypeList() []QueryMediaContentType {
	return slices.Clone(contentTypeList)
}
