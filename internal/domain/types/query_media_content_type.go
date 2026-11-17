package dtypes

import (
	"errors"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/neosy/elengrab/internal/pkg/stringx"
)

type QueryMediaContentType string

const (
	QueryMediaContentTypeNone   QueryMediaContentType = ""
	QueryMediaContentTypeAll    QueryMediaContentType = "all"
	QueryMediaContentTypeVideo  QueryMediaContentType = "video"
	QueryMediaContentTypeShorts QueryMediaContentType = "shorts"
	QueryMediaContentTypeAudio  QueryMediaContentType = "audio"
)

var (
	queryMediaContentTypes = map[QueryMediaContentType]struct{}{
		QueryMediaContentTypeAll:    {},
		QueryMediaContentTypeVideo:  {},
		QueryMediaContentTypeShorts: {},
		QueryMediaContentTypeAudio:  {},
	}

	queryMediaContentTypeByName = map[string]QueryMediaContentType{
		"all":    QueryMediaContentTypeAll,
		"video":  QueryMediaContentTypeVideo,
		"shorts": QueryMediaContentTypeShorts,
		"audio":  QueryMediaContentTypeAudio,
	}

	contentTypeList = []QueryMediaContentType{
		QueryMediaContentTypeAll,
		QueryMediaContentTypeVideo,
		QueryMediaContentTypeShorts,
		QueryMediaContentTypeAudio,
	}
)

// String returns the value as a string.
func (v QueryMediaContentType) String() string {
	return string(v)
}

// Label returns the display label for the content type.
func (v QueryMediaContentType) Label() string {
	return stringx.Capitalize(string(v))
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
