package dtypes

import (
	"errors"
	"strings"
)

type QueryFieldName string

const (
	QueryFieldNameNone QueryFieldName = ""

	QueryFieldNameUserID    QueryFieldName = "userID"
	QueryFieldNameTitle     QueryFieldName = "title"
	QueryFieldNameCreatedAt QueryFieldName = "createdAt"
)

var (
	parseQueryFieldNameMap = map[string]QueryFieldName{
		"userID":    QueryFieldNameUserID,
		"title":     QueryFieldNameTitle,
		"createdAt": QueryFieldNameCreatedAt,
	}
)

// String returns the string representation of the QueryFieldName.
func (name QueryFieldName) String() string {
	return string(name)
}

// ParseQueryFieldName parses a string into a QueryFieldName.
// Returns an error if the string does not correspond to a valid QueryFieldName.
func ParseQueryFieldName(s string) (QueryFieldName, error) {
	fieldName, exists := parseQueryFieldNameMap[strings.ToLower(s)]
	if !exists {
		return QueryFieldNameNone, errors.New("invalid value for QueryFieldName")
	}
	return fieldName, nil
}
