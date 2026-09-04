package dbutils

import (
	"reflect"
	"time"

	"github.com/Masterminds/squirrel"
)

// DateFilter is a struct that represents a filter for a date field in a database query.
type DateFilter struct {
	// Date is the date to filter by.
	Date time.Time
	// Direction is the direction of the filter.
	Direction FilterOperator
}

// NewDateFilter creates a new DateFilter with the given date and direction.
func NewDateFilter(date time.Time, direction FilterOperator) FilterConditioner {
	return &DateFilter{
		Date:      date,
		Direction: direction,
	}
}

// Value returns the value of the comparison filter.
func (f *DateFilter) Value() any {
	return f.Date
}

// Operator returns the operator of the comparison filter.
func (f *DateFilter) Operator() FilterOperator {
	return f.Direction
}

// SqlCondition returns a squirrel.Condition for the date filter
func (f *DateFilter) SqlCondition(fieldName string) squirrel.Sqlizer {
	return NewFilterCondition(f.Date, f.Direction).SqlCondition(fieldName)
}

// IsZero returns true if the date filter is zero, false otherwise.
func (f *DateFilter) IsZero() bool {
	if f == nil {
		return true
	}
	return reflect.ValueOf(f.Date).IsZero()
}

// IsSet returns true if the date filter is set, false otherwise.
func (f *DateFilter) IsSet() bool {
	return !f.IsZero()
}

// Condition returns a FilterConditioner for the date filter.
func (f *DateFilter) Condition() FilterConditioner {
	if f == nil {
		return nil
	}
	return NewFilterCondition(f.Date, f.Direction)
}
