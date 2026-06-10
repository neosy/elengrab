package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Application mode
type AppMode uint8

const (
	// Public: anonymous access allowed, downloads visible to everyone.
	AppModePublic AppMode = iota

	// Guest: anonymous access allowed, downloads visible only to their owner.
	AppModeGuest

	// Authenticated: login required, access controlled by user permissions.
	AppModeAuthenticated

	AppModeDefault = AppModePublic
)

var (
	// appModeMap implementation of a set for AppMode
	appModeMap = map[AppMode]string{
		AppModePublic:        "public",
		AppModeGuest:         "guest",
		AppModeAuthenticated: "authenticated",
	}

	parseAppModeMap = map[string]AppMode{
		"public":        AppModePublic,
		"guest":         AppModeGuest,
		"authenticated": AppModeAuthenticated,
	}

	appModeToUniquenessScopeMap = map[AppMode]UniquenessScope{
		AppModePublic:        UniquenessScopeGlobal,
		AppModeGuest:         UniquenessScopePerUser,
		AppModeAuthenticated: UniquenessScopePerUser,
	}
)

// String returns the value as a string.
func (v AppMode) String() string {
	return appModeMap[v]
}

// Exists returns true if the AppMode is valid.
func (v AppMode) Exists() bool {
	_, exists := appModeMap[v]
	return exists
}

// UniquenessScope returns the corresponding UniquenessScope for the AppMode.
func (v AppMode) UniquenessScope() UniquenessScope {
	return appModeToUniquenessScopeMap[v]
}

// IsGuestAllowed returns true if guests are allowed for this app mode.
func (v AppMode) IsGuestAllowed() bool {
	return v != AppModeAuthenticated
}

// IsUserRequired returns true if user is required for this app mode.
func (v AppMode) IsUserRequired() bool {
	return v == AppModeAuthenticated || v == AppModeGuest
}

// ParseAppMode converting string to AppMode
func ParseAppMode(s string) (AppMode, error) {
	mode, exists := parseAppModeMap[strings.ToLower(s)]
	if !exists {
		return AppModeDefault, errors.New("invalid value for AppMode")
	}
	return mode, nil
}

// MustParseAppMode converting string to AppMode, ignoring any errors.
func MustParseAppMode(s string) AppMode {
	mode, err := ParseAppMode(s)
	if err != nil {
		return AppModeDefault
	}
	return mode
}

// ValidateAppMode checks if the field value is a valid AppMode enum.
func ValidateAppMode(fl validator.FieldLevel) bool {
	_, exist := parseAppModeMap[fl.Field().String()]
	return exist
}
