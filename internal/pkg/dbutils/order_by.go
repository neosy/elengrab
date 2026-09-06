package dbutils

import (
	"slices"
	"strings"
)

const (
	orderSep = ", "
)

// OrderBy represents a single sorting criterion for SQL queries,
// specifying the field to sort by and the direction of sorting (ascending or descending).
type OrderBy struct {
	Alias     string
	Field     string
	Direction OrderDirection
}

// OrderByList is a list of OrderBy structs representing multiple sorting criteria.
type OrderByList []OrderBy

// SortBy creates an OrderBy struct for the specified field and order direction.
func SortBy(field string, direction OrderDirection) OrderBy {
	return OrderBy{
		Field:     field,
		Direction: direction,
	}
}

// SortByWithAlias creates an OrderBy struct for the specified field, order direction, and table alias.
func SortByWithAlias(field string, direction OrderDirection, alias string) OrderBy {
	return OrderBy{
		Alias:     alias,
		Field:     field,
		Direction: direction,
	}
}

// QueryOrderByList constructs a list of SQL ORDER BY clauses from a list of OrderBy structs.
// It returns a slice of strings, each representing a sorting criterion in the format "field direction".
// If the list is empty, it returns an empty slice.
func QueryOrderByList(list ...OrderBy) []string {
	queryOrderList := make([]string, 0, len(list))

	for _, orderBy := range list {
		query := orderBy.Field + " " + orderBy.Direction.String()
		if orderBy.Alias != "" {
			query = orderBy.Alias + "." + query
		}
		queryOrderList = append(queryOrderList, query)
	}

	return queryOrderList
}

// QueryOrderBy constructs a SQL ORDER BY clause from a list of OrderBy structs.
// It returns a string that can be used in SQL queries to specify the sorting order.
// If the list is empty, it returns an empty string.
func QueryOrderBy(list ...OrderBy) string {
	return strings.Join(QueryOrderByList(list...), orderSep)
}

// List returns an OrderByList containing the current OrderBy struct.
func (ob OrderBy) List() OrderByList {
	return OrderByList{ob}
}

// Query constructs a SQL ORDER BY clause from a single OrderBy struct.
// It returns a string that can be used in SQL queries to specify the sorting order.
func (ob OrderBy) Query() string {
	return QueryOrderBy(ob)
}

// WithAlias returns a new OrderBy struct with the specified alias.
func (ob OrderBy) WithAlias(alias string) OrderBy {
	ob.Alias = alias
	return ob
}

// Query constructs a SQL ORDER BY clause from an OrderByList.
// It returns a string that can be used in SQL queries to specify the sorting order.
// If the list is empty, it returns an empty string.
func (l OrderByList) Query() string {
	return QueryOrderBy(l...)
}

// WithAlias returns a new OrderByList with the specified alias applied to each OrderBy struct in the list.
func (l OrderByList) WithAlias(alias string) OrderByList {
	copyList := slices.Clone(l)

	for i := range copyList {
		copyList[i].Alias = alias
	}
	return copyList
}

// QueryList constructs a list of SQL ORDER BY clauses from an OrderByList.
// It returns a slice of strings, each representing a sorting criterion in the format "field direction".
// If the list is empty, it returns an empty slice.
func (l OrderByList) QueryList() []string {
	return QueryOrderByList(l...)
}

// Append appends an OrderBy struct to the OrderByList.
func (l *OrderByList) Append(orderBy OrderBy) {
	if l == nil {
		return
	}
	*l = append(*l, orderBy)
}

// Add adds a new OrderBy struct to the OrderByList with the specified field and order direction.
func (l *OrderByList) Add(field string, direction OrderDirection) {
	if l == nil {
		return
	}

	orderBy := OrderBy{
		Field:     field,
		Direction: direction,
	}

	l.Append(orderBy)
}
