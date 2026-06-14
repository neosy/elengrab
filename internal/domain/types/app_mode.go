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

	// AppModePublicReadonly: anonymous access allowed, but only read-only for public media.
	// Uploading/adding new content requires authentication.
	AppModePublicReadonly

	// Authenticated: login required, access controlled by user permissions.
	AppModeAuthenticated

	AppModeDefault = AppModePublic
)

var (
	// nameByAppMode maps an AppMode to its string representation
	nameByAppMode = map[AppMode]string{
		AppModePublic:         "public",
		AppModeGuest:          "guest",
		AppModePublicReadonly: "public_readonly",
		AppModeAuthenticated:  "authenticated",
	}

	appModeByName = map[string]AppMode{
		"public":          AppModePublic,
		"guest":           AppModeGuest,
		"public_readonly": AppModePublicReadonly,
		"authenticated":   AppModeAuthenticated,
	}

	// uniquenessScopeByAppMode maps each AppMode to its corresponding uniqueness scope.
	uniquenessScopeByAppMode = map[AppMode]UniquenessScope{
		AppModePublic:         UniquenessScopeGlobal,
		AppModeGuest:          UniquenessScopePerUser,
		AppModePublicReadonly: UniquenessScopePerUser,
		AppModeAuthenticated:  UniquenessScopePerUser,
	}
)

// String returns the value as a string.
func (v AppMode) String() string {
	return nameByAppMode[v]
}

// Exists returns true if the AppMode is valid.
func (v AppMode) Exists() bool {
	_, exists := nameByAppMode[v]
	return exists
}

// UniquenessScope returns the corresponding UniquenessScope for the AppMode.
func (v AppMode) UniquenessScope() UniquenessScope {
	return uniquenessScopeByAppMode[v]
}

// IsGuestAllowed returns true if guests are allowed for this app mode.
func (v AppMode) IsGuestAllowed() bool {
	return v == AppModePublic || v == AppModeGuest
}

// IsUserRequiredForWrite1 returns true if user is required for this app mode.
func (v AppMode) IsUserRequiredForWrite() bool {
	return v == AppModeAuthenticated || v == AppModeGuest || v == AppModePublicReadonly
}

// IsUserRequiredForRead returns true if user is required for this app mode.
func (v AppMode) IsUserRequiredForRead() bool {
	return v == AppModeAuthenticated
}

// ParseAppMode converting string to AppMode
func ParseAppMode(s string) (AppMode, error) {
	mode, exists := appModeByName[strings.ToLower(s)]
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
	_, exist := appModeByName[fl.Field().String()]
	return exist
}
