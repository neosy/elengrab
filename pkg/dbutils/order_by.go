package dbutils

import "strings"

const (
	orderSep = ", "

	OrderAsc  = "ASC"
	OrderDesc = "DESC"
)

// Flds is a map where the key is a field name (string), and the value is the associated parameter (string).
// It is typically used to store field names and their respective parameters for queries or configuration.
type Flds = map[string]string

// OrderBy creates a string for sorting with multiple fields and directions.
// The function accepts a map where keys are the field names and values are the sorting directions ("ASC" or "DESC").
// It returns a comma-separated string that can be used in the ORDER BY clause of a SQL query.
//
// Example usage:
// pgutils.OrderBy(
//
//	pgutils.Flds{"sales_id": pgutils.OrderDesc},  // Sort by sales_id in ascending sales
//	pgutils.Flds{"created_at": pgutils.OrderAsc}, // Sort by created_at in descending sales
//
// )
func OrderBy(fields ...Flds) string {
	var orderFields []string
	for _, fieldDirs := range fields {
		for fieldName, direction := range fieldDirs {
			orderFields = append(orderFields, fieldName+" "+direction)
		}
	}
	return strings.Join(orderFields, orderSep)
}
