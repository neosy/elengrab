package dbutils

import (
	"reflect"

	"github.com/Masterminds/squirrel"
)

// FilterConditioner is an interface that defines methods for a filter condition.
type FilterConditioner interface {
	Value() any
	Operator() FilterOperator
	IsZero() bool
	IsSet() bool
	SqlCondition(fieldName string) squirrel.Sqlizer
}

// FilterCondition represents a filter that compares a value using an operator.
type FilterCondition[T any] struct {
	// value is the value to compare against.
	value T
	// operator is the comparison operator to use (e.g., less than, greater than, equal).
	operator FilterOperator
}

// NewFilterCondition creates a new FilterCondition with the given value and operator.
func NewFilterCondition[T any](value T, operator FilterOperator) FilterConditioner {
	return &FilterCondition[T]{
		value:    value,
		operator: operator,
	}
}

// NewEqFilterCondition creates a new FilterCondition with the Eq operator.
func NewFilterConditionEq[T any](value T) FilterConditioner {
	return NewFilterCondition(value, FilterOperatorEq)
}

// Value returns the value of the comparison filter.
func (f *FilterCondition[T]) Value() any {
	if f == nil {
		return nil
	}
	return f.value
}

// Operator returns the operator of the comparison filter.
func (f *FilterCondition[T]) Operator() FilterOperator {
	if f == nil {
		return FilterOperatorEq
	}
	return f.operator
}

// SqlCondition generates a SQL condition for the comparison filter.
// It returns a squirrel.Sqlizer that can be used in SQL queries.
// The fieldName parameter specifies the name of the database field to compare against.
func (f FilterCondition[T]) SqlCondition(fieldName string) squirrel.Sqlizer {
	switch f.operator {
	case FilterOperatorEq:
		return squirrel.Eq{fieldName: f.value}
	case FilterOperatorNotEq:
		return squirrel.NotEq{fieldName: f.value}
	case FilterOperatorLt:
		return squirrel.Lt{fieldName: f.value}
	case FilterOperatorLtOrEq:
		return squirrel.LtOrEq{fieldName: f.value}
	case FilterOperatorGt:
		return squirrel.Gt{fieldName: f.value}
	case FilterOperatorGtOrEq:
		return squirrel.GtOrEq{fieldName: f.value}
	case FilterOperatorLike:
		return squirrel.Like{fieldName: f.value}
	default:
		return squirrel.Eq{fieldName: f.value}
	}
}

// IsZero checks if the comparison filter is zero (i.e., has a zero value).
// It returns true if the filter is zero, false otherwise.
func (f *FilterCondition[T]) IsZero() bool {
	if f == nil {
		return true
	}
	return reflect.ValueOf(f.value).IsZero()
}

// IsSet checks if the comparison filter is set (i.e., has a non-zero value).
// It returns true if the filter is set, false otherwise.
func (f *FilterCondition[T]) IsSet() bool {
	return !f.IsZero()
}
