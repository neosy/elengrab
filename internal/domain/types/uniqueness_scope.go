package dtypes

import (
	"errors"
	"strings"
)

type UniquenessScope uint8

const (
	// Unique across the entire system
	UniquenessScopeGlobal UniquenessScope = iota
	// Unique within a single user scope
	UniquenessScopePerUser

	UniquenessScopeDefault = UniquenessScopeGlobal
)

var uniquenessScopeMap = map[UniquenessScope]string{
	UniquenessScopeGlobal:  "global",
	UniquenessScopePerUser: "per_user",
}

var parseUniquenessScopeMap = map[string]UniquenessScope{
	"global":   UniquenessScopeGlobal,
	"per_user": UniquenessScopePerUser,
}

// String returns the value as a string.
func (v UniquenessScope) String() string {
	return uniquenessScopeMap[v]
}

// Exists returns true if the UniquenessScope is valid.
func (v UniquenessScope) Exists() bool {
	_, exists := uniquenessScopeMap[v]
	return exists
}

// ParseUniquenessScope converting string to UniquenessScope
func ParseUniquenessScope(s string) (UniquenessScope, error) {
	scope, exists := parseUniquenessScopeMap[strings.ToLower(s)]
	if !exists {
		return UniquenessScopeDefault, errors.New("invalid value for HistoryMode")
	}
	return scope, nil
}

// MustParseUniquenessScope converting string to UniquenessScope, ignoring any errors.
func MustParseUniquenessScope(s string) UniquenessScope {
	scope, _ := ParseUniquenessScope(s)
	return scope
}
