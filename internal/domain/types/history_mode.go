package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// History mode
type HistoryMode string

const (
	HistoryModeDisabled HistoryMode = "disabled"
	HistoryModeGlobal   HistoryMode = "global"
	HistoryModePerUser  HistoryMode = "per_user"

	HistoryModeDefault = HistoryModeGlobal
)

var (
	// historyModeMap implementation of a set for HistoryMode
	historyModeMap = map[HistoryMode]struct{}{
		HistoryModeDisabled: {},
		HistoryModeGlobal:   {},
		HistoryModePerUser:  {},
	}
)

// String returns the value as a string.
func (v HistoryMode) String() string {
	return string(v)
}

// Exists returns true if the HistoryMode is valid.
func (v HistoryMode) Exists() bool {
	_, exists := historyModeMap[v]
	return exists
}

// ParseHistoryMode converting string to status
func ParseHistoryMode(s string) (HistoryMode, error) {
	mode := HistoryMode(strings.ToLower(s))

	if _, exists := historyModeMap[mode]; !exists {
		return "", errors.New("invalid value for HistoryMode")
	}

	return mode, nil
}

func MustParseHistoryMode(s string) HistoryMode {
	mode, err := ParseHistoryMode(s)
	if err != nil {
		return HistoryModeDefault
	}

	return mode
}

// StringToHistoryMode converting string to status
func StringToHistoryMode(s string) (HistoryMode, error) {
	return ParseHistoryMode(s)
}

// ValidateHistoryMode checks if the field value is a valid HistoryMode enum.
func ValidateHistoryMode(fl validator.FieldLevel) bool {
	return HistoryMode(fl.Field().String()).Exists()
}
