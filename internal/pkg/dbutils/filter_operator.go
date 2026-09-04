package dbutils

// FilterOperator represents a comparison operator used by a filter.
type FilterOperator uint8

const (
	// FilterOperatorEq represents the equal operator (=).
	FilterOperatorEq FilterOperator = iota

	// FilterOperatorNotEq represents the not-equal operator (!=).
	FilterOperatorNotEq

	// FilterOperatorLt represents the less-than operator (<).
	FilterOperatorLt

	// FilterOperatorLtOrEqual represents the less-than-or-equal operator (<=).
	FilterOperatorLtOrEq

	// FilterOperatorGt represents the greater-than operator (>).
	FilterOperatorGt

	// FilterOperatorGtOrEqual represents the greater-than-or-equal operator (>=).
	FilterOperatorGtOrEq

	//FilterOperatorGtOrLike represents the LIKE operator for pattern matching.
	FilterOperatorLike
)
