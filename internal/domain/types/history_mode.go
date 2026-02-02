package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// History mode
type HistoryMode uint8

const (
	// No history is shown
	HistoryModeDisabled HistoryMode = iota
	// History is shared globally
	HistoryModeGlobal
	// History is unique per user
	HistoryModePerUser

	HistoryModeDefault = HistoryModeGlobal
)

var (
	// historyModeMap implementation of a set for HistoryMode
	historyModeMap = map[HistoryMode]string{
		HistoryModeDisabled: "disabled",
		HistoryModeGlobal:   "global",
		HistoryModePerUser:  "per_user",
	}

	parseHistoryModeMap = map[string]HistoryMode{
		"disabled": HistoryModeDisabled,
		"global":   HistoryModeGlobal,
		"per_user": HistoryModePerUser,
	}

	historyModeToUniquenessScopeMap = map[HistoryMode]UniquenessScope{
		HistoryModeDisabled: UniquenessScopeGlobal,
		HistoryModeGlobal:   UniquenessScopeGlobal,
		HistoryModePerUser:  UniquenessScopePerUser,
	}
)

// String returns the value as a string.
func (v HistoryMode) String() string {
	return historyModeMap[v]
}

// Exists returns true if the HistoryMode is valid.
func (v HistoryMode) Exists() bool {
	_, exists := historyModeMap[v]
	return exists
}

// UniquenessScope returns the corresponding UniquenessScope for the HistoryMode.
func (v HistoryMode) UniquenessScope() UniquenessScope {
	return historyModeToUniquenessScopeMap[v]
}

// ParseHistoryMode converting string to HistoryMode
func ParseHistoryMode(s string) (HistoryMode, error) {
	mode, exists := parseHistoryModeMap[strings.ToLower(s)]
	if !exists {
		return HistoryModeDefault, errors.New("invalid value for HistoryMode")
	}
	return mode, nil
}

// MustParseHistoryMode converting string to HistoryMode, ignoring any errors.
func MustParseHistoryMode(s string) HistoryMode {
	mode, err := ParseHistoryMode(s)
	if err != nil {
		return HistoryModeDefault
	}
	return mode
}

// ValidateHistoryMode checks if the field value is a valid HistoryMode enum.
func ValidateHistoryMode(fl validator.FieldLevel) bool {
	_, exist := parseHistoryModeMap[fl.Field().String()]
	return exist
}
