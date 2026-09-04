package dtypes

import "errors"

type QueryFilterName string

const (
	QueryFilterNameNone QueryFilterName = ""

	QueryFilterNameUserID    QueryFilterName = "userID"
	QueryFilterNameTitle     QueryFilterName = "title"
	QueryFilterNameCreatedAt QueryFilterName = "createdAt"

	QueryFilterNameSearch      QueryFilterName = "search"
	QueryFilterNameDownloadIDs QueryFilterName = "downloadIDs"
)

var (
	parseQueryFilterNameMap = map[string]QueryFilterName{
		"userID":      QueryFilterNameUserID,
		"title":       QueryFilterNameCreatedAt,
		"createdAt":   QueryFilterNameTitle,
		"search":      QueryFilterNameSearch,
		"downloadIDs": QueryFilterNameDownloadIDs,
	}
)

// String returns the string representation of the QueryFilterName.
func (name QueryFilterName) String() string {
	return string(name)
}

// ParseQueryFilterName parses a string into a QueryFilterName.
// Returns an error if the string does not correspond to a valid QueryFilterName.
func ParseQueryFilterName(s string) (QueryFilterName, error) {
	filterName, exists := parseQueryFilterNameMap[s]
	if !exists {
		return QueryFilterNameNone, errors.New("invalid value for QueryFilterName")
	}
	return filterName, nil
}
