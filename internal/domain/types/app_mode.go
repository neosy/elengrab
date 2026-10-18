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

// AllowsGuests reports whether guest access is permitted in this mode.
func (m AppMode) AllowsGuests() bool {
	return m == AppModePublic || m == AppModeGuest
}

// RequiresIdentityForWrite reports whether write operations need an identified client
// (guest or authenticated user) in this mode.
func (m AppMode) RequiresIdentityForWrite() bool {
	return m == AppModeAuthenticated || m == AppModeGuest || m == AppModePublicReadonly
}

// RequiresIdentityForRead reports whether read operations need an identified client
// (guest or authenticated user) in this mode.
func (m AppMode) RequiresIdentityForRead() bool {
	return m == AppModeAuthenticated
}

// AllowsAnonymousWrite reports whether unidentified clients may write in this mode.
// It is the inverse of RequiresIdentityForWrite.
func (m AppMode) AllowsAnonymousWrite() bool {
	return !m.RequiresIdentityForWrite()
}

// AllowsAnonymousRead reports whether unidentified clients may read in this mode.
// It is the inverse of RequiresIdentityForRead.
func (m AppMode) AllowsAnonymousRead() bool {
	return !m.RequiresIdentityForRead()
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
