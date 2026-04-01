package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Application mode
type AppMode uint8

const (
	// Public: guests allowed, data visible to everyone
	AppModePublic AppMode = iota
	// PerUser: guests allowed, data visible only to the owner
	AppModePerUser
	// AuthOnly: guests forbidden, login required
	AppModeAuthOnly

	AppModeDefault = AppModePublic
)

var (
	// appModeMap implementation of a set for AppMode
	appModeMap = map[AppMode]string{
		AppModePublic:   "public",
		AppModePerUser:  "per_user",
		AppModeAuthOnly: "auth_only",
	}

	parseAppModeMap = map[string]AppMode{
		"public":    AppModePublic,
		"per_user":  AppModePerUser,
		"auth_only": AppModeAuthOnly,
	}

	appModeToUniquenessScopeMap = map[AppMode]UniquenessScope{
		AppModePublic:   UniquenessScopeGlobal,
		AppModePerUser:  UniquenessScopePerUser,
		AppModeAuthOnly: UniquenessScopePerUser,
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

// ParseAppMode converting string to AppMode
func ParseAppMode(s string) (AppMode, error) {
	mode, exists := parseAppModeMap[strings.ToLower(s)]
	if !exists {
		return AppModeDefault, errors.New("invalid value for AppMode")
	}
	return mode, nil
}

// IsGuestAllowed returns true if guests are allowed for this app mode.
func (v AppMode) IsGuestAllowed() bool {
	return v != AppModeAuthOnly
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
