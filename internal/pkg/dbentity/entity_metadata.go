package dbentity

import (
	"maps"
	"slices"
)

// EntityMetadata holds metadata about a database entity,
// including its query fields, insert fields, searchable fields, and pagination field names.
type EntityMetadata struct {
	// queryFields holds the list of fields that will be used for queries.
	queryFields []string
	// insertFields holds the list of fields that will be included in insert operations.
	insertFields []string
	// searchableFields holds the list of fields explicitly marked as searchable.
	searchableFields []string
	// paginationFieldNames holds a map of field names to their corresponding pagination field names.
	paginationFieldNames map[string]string
}

// NewEntityMetadata creates a new EntityMetadata instance based on the provided BaseEntity.
func NewEntityMetadata[T any](e BaseEntity[T]) EntityMetadata {
	return EntityMetadata{
		queryFields:          e.QueryFields(),
		insertFields:         e.InsertFields(),
		searchableFields:     e.SearchableFields(),
		paginationFieldNames: e.PaginationFieldNames(),
	}
}

// QueryFields returns a list of fields that will be used for queries
func (m *EntityMetadata) QueryFields() []string {
	return slices.Clone(m.queryFields)
}

// InsertFields returns fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (m *EntityMetadata) InsertFields() []string {
	return slices.Clone(m.insertFields)
}

// SearchableFields returns a list of fields explicitly marked as searchable
func (m *EntityMetadata) SearchableFields() []string {
	return slices.Clone(m.searchableFields)
}

// PaginationFieldNames returns a map of field names to their corresponding pagination field names.
func (m *EntityMetadata) PaginationFieldNames() map[string]string {
	return maps.Clone(m.paginationFieldNames)
}
