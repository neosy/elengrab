package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type QueryMediaViewMode string

const (
	QueryMediaViewModeNone    QueryMediaViewMode = ""
	QueryMediaViewModeNew     QueryMediaViewMode = "new"
	QueryMediaViewModePopular QueryMediaViewMode = "popular"
	QueryMediaViewModeOld     QueryMediaViewMode = "old"

	QueryMediaViewModeDefault = QueryMediaViewModeNew
)

var (
	queryMediaViewModes = map[QueryMediaViewMode]struct{}{
		QueryMediaViewModeNew:     {},
		QueryMediaViewModePopular: {},
		QueryMediaViewModeOld:     {},
	}

	queryMediaViewModeByName = map[string]QueryMediaViewMode{
		"new":     QueryMediaViewModeNew,
		"popular": QueryMediaViewModePopular,
		"old":     QueryMediaViewModeOld,
	}
)

// String returns the value as a string.
func (v QueryMediaViewMode) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v QueryMediaViewMode) Ptr() *QueryMediaViewMode {
	return &v
}

// Exists returns true if the QueryMediaViewMode is valid.
func (v QueryMediaViewMode) Exists() bool {
	_, exists := queryMediaViewModes[v]
	return exists
}

// ParseQueryMediaViewMode converting string to QueryMediaViewMode
func ParseQueryMediaViewMode(s string) (QueryMediaViewMode, error) {
	queryMediaViewMode, exists := queryMediaViewModeByName[strings.ToLower(s)]
	if !exists {
		return QueryMediaViewModeDefault, errors.New("invalid value for QueryMediaViewMode")
	}

	return queryMediaViewMode, nil
}

// ValidateQueryMediaViewMode checks if the field value is a valid QueryMediaViewMode enum.
func ValidateQueryMediaViewMode(fl validator.FieldLevel) bool {
	_, err := ParseQueryMediaViewMode(fl.Field().String())
	return err == nil
}
