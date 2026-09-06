package dbutils

// FieldFilter represents a filter condition for a specific field in a database query.
type FieldFilter struct {
	// Field is the name of the field to filter on.
	Field string

	// Value is the value to filter by for the specified field.
	Value any
}

// FiltersByField is a map that associates field names with their corresponding FieldFilter.
type FiltersByField map[string]FieldFilter

// Add adds a new FieldFilter to the FiltersByField map for the specified field and value.
func (filters FiltersByField) Add(field string, value any) {
	if filters == nil {
		return
	}

	filter := FieldFilter{
		Field: field,
		Value: value,
	}

	filters[field] = filter
}
