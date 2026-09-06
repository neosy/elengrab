package dbutils

// OrderDirection represents the direction of sorting for SQL queries, either ascending or descending.
type OrderDirection uint8

const (
	// OrderAscending represents ascending order for sorting.
	// It is used in SQL queries to specify that results should be sorted from lowest to highest.
	// Example: ORDER BY column_name ASC
	OrderAscending OrderDirection = iota

	// OrderDescending represents descending order for sorting.
	// It is used in SQL queries to specify that results should be sorted from highest to lowest.
	// Example: ORDER BY column_name DESC
	OrderDescending
)

var (
	// orderDirectionToString maps OrderDirection values to their corresponding string representations.
	orderDirectionToString = map[OrderDirection]string{
		OrderAscending:  "ASC",
		OrderDescending: "DESC",
	}
)

// String returns the string representation of the OrderDirection.
func (o OrderDirection) String() string {
	return orderDirectionToString[o]
}
