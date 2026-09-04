package dbutils

import (
	"github.com/Masterminds/squirrel"
)

// FieldFilter represents a filter condition for a specific field in a database query.
type FieldFilter struct {
	// FieldName is the name of the field to filter on.
	FieldName string

	// Condition is the filter condition to apply to the field.
	condition FilterConditioner
}

// FiltersByField is a map that associates field names with their corresponding FieldFilter.
type FiltersByField map[string]FieldFilter

// NewFieldFilter creates a new FieldFilter with the specified field name, value, and operator.
func NewFieldFilter(fieldName string, value any, operator FilterOperator) FieldFilter {
	return FieldFilter{
		FieldName:     fieldName,
		condition: NewFilterCondition(value, operator),
	}
}

// NewFieldFilterEq creates a new FieldFilter with the Eq operator.
func NewFieldFilterEq(fieldName string, value any) FieldFilter {
	return NewFieldFilter(fieldName, value, FilterOperatorEq)
}

func (f FieldFilter) Condition() FilterConditioner {
	return f.condition
}

// Append adds a new FieldFilter to the FiltersByField map for the specified field and condition.
func (filters FiltersByField) Append(filter FieldFilter) FiltersByField {
	if filters == nil {
		filters = make(FiltersByField)
	}

	filters[filter.FieldName] = filter
	return filters
}

// Add adds a new FieldFilter to the FiltersByField map for the specified field and condition.
func (filters FiltersByField) Add(field string, condition FilterConditioner) {
	if filters == nil {
		return
	}

	if condition == nil {
		condition = NewFilterConditionEq[any](nil)
	}

	filter := FieldFilter{
		FieldName:     field,
		condition: condition,
	}

	filters[field] = filter
}

// SqlCondition generates a SQL condition for the FieldFilter using the provided field name and condition.
func (f FieldFilter) SqlCondition() squirrel.Sqlizer {
	if f.condition == nil {
		return squirrel.Eq{f.FieldName: nil}
	}
	return f.condition.SqlCondition(f.FieldName)
}

// SqlConditionWithAlias generates a SQL condition for the FieldFilter using the provided field name, alias, and condition.
func (f FieldFilter) SqlConditionWithAlias(alias string) squirrel.Sqlizer {
	name := f.FieldName

	if alias != "" {
		name = alias + "." + name
	}

	if f.condition == nil {
		return squirrel.Eq{name: nil}
	}

	return f.condition.SqlCondition(name)
}
