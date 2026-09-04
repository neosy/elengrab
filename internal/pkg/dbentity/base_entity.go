package dbentity

import (
	etags "github.com/neosy/elengrab/internal/pkg/dbentity/tags"
	"github.com/neosy/elengrab/internal/pkg/reflection"
)

// BaseEntity is a base structure that provides common methods
// for working with entity fields and their SQL tag names
type BaseEntity[T any] struct{}

// FieldNameFromPointer field name from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId)
func (e *BaseEntity[T]) FieldNameFromPointer(structPtr *T, fieldPtr any) string {
	return e.FieldNameFromPointerWithAlias(structPtr, fieldPtr, "")
}

// FieldNameFromPointerWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *BaseEntity[T]) FieldNameFromPointerWithAlias(structPtr *T, fieldPtr any, alias string) string {
	name, _ := reflection.StructFieldName(structPtr, fieldPtr, etags.ColumnTagName().String())
	return e.FieldNameWithAlias(name, alias)
}

// FieldName returns the field name from the SQL tag using a structure field name or pointer,
// optionally prefixed with a table alias.
//
// Examples:
//
//	var entity <TableEntity>
//	entity.FieldName("created_at", "")              // "created_at"
//	entity.FieldName(&entity.CreatedAt, "")         // "created_at"
//	entity.FieldName(&entity.CreatedAt, "users")    // "users.created_at"
func (e *BaseEntity[T]) FieldName(structPtr *T, field any, alias ...string) string {
	var fieldName string

	switch v := field.(type) {
	case string:
		fieldName = v
	default:
		fieldName = e.FieldNameFromPointer(structPtr, v)
	}

	if len(alias) == 0 || alias[0] == "" {
		return fieldName
	}

	return alias[0] + "." + fieldName
}

// FieldNameWithAlias returns the field name with alias
//
// Examples:
//
//	 var ent <TableEntity>
//		ent.FieldNameWithAlias("created_at", "")       // "created_at"
//		ent.FieldNameWithAlias("created_at", "users")   // "users.created_at"
func (e *BaseEntity[T]) FieldNameWithAlias(fieldName, alias string) string {
	if alias == "" {
		return fieldName
	}
	return alias + "." + fieldName
}

// PaginateFieldNameFromPointer returns the field name with alias from the pagination tag (`pfield`)
// based on the pointer to the struct and the pointer to the field.
// Example:
// var ent <TableEntity>
// ent.PaginateFieldName(&ent, &ent.SalesId)
func (e *BaseEntity[T]) PaginateFieldNameFromPointer(structPtr *T, fieldPtr any) string {
	name, _ := reflection.StructFieldName(structPtr, fieldPtr, etags.TagNamePaginationField.String())
	return name
}

// QueryFields returns a list of fields that will be used for queries
func (e *BaseEntity[T]) QueryFields() []string {
	var ent T
	return etags.FieldNames(&ent, etags.TagNameSelect)
}

// SearchableFields returns a list of fields explicitly marked as searchable
func (e *BaseEntity[T]) SearchableFields() []string {
	var ent T
	return etags.FieldsWithTrueTag(&ent, etags.TagNameIsSearch)
}

// FieldsWithAlias returns a list of fields with alias
func (e *BaseEntity[T]) FieldsWithAlias(fields []string, alias string) []string {
	withAlias := make([]string, len(fields))
	for i, field := range fields {
		withAlias[i] = alias + "." + field
	}

	return withAlias
}

// QueryFieldsWithAlias returns a list of fields with alias that will be used for queries
func (e *BaseEntity[T]) QueryFieldsWithAlias(alias string) []string {
	return e.FieldsWithAlias(e.QueryFields(), alias)
}

// InsertFields returns fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *BaseEntity[T]) InsertFields() []string {
	var ent T
	return etags.FieldNamesExceptFalseTag(&ent, etags.TagNameInsert)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *BaseEntity[T]) InsertValues(structPtr *T) []any {
	return etags.ValuesExceptFalseTag(structPtr, etags.TagNameInsert)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
// This can be used for scanning database rows into the struct.
func (e *BaseEntity[T]) FieldPointers(structPtr *T) ([]any, error) {
	return reflection.StructFieldPointers(structPtr)
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *BaseEntity[T]) FieldPointer(structPtr *T, fieldName string) (any, error) {
	return reflection.StructFieldPointer(structPtr, fieldName, etags.ColumnTagName().String())
}

// FieldValues returns a map of field names to their corresponding values
func (e *BaseEntity[T]) FieldValues(fields []string, values []any) map[string]any {
	m := make(map[string]any, len(fields))

	for i, f := range fields {
		m[f] = values[i]
	}

	return m
}

// InsertFieldValues returns a map of field names to their corresponding values
// for fields included in insert operations.
func (e *BaseEntity[T]) InsertFieldValues(structPtr *T) map[string]any {
	return e.FieldValues(e.InsertFields(), e.InsertValues(structPtr))
}

// FieldNamesByTag returns a map that maps field names from the specified tag
// to their corresponding database column names.
//
// For example, given the field:
//
//	SourceCreatedAt time.Time `db:"created_at" pfield:"createdAt"`
//
// calling FieldNamesByTag with the "pfield" tag returns:
//
//	map[string]string{"createdAt": "created_at"}
func (e *BaseEntity[T]) FieldNamesByTag(tag etags.TagName) map[string]string {
	var ent T
	return etags.FieldNamesByTags(&ent, tag, etags.ColumnTagName())
}

// PaginationFieldNames returns a map from pagination field names to their corresponding database column names.
// It uses the `pfield` tag to determine the pagination field names.
//
// For example, given the field:
//
//	SourceCreatedAt time.Time `db:"created_at" pfield:"createdAt"`
//
// calling PaginationFieldNames returns:
//
//	map[string]string{"createdAt": "created_at"}
func (e *BaseEntity[T]) PaginationFieldNames() map[string]string {
	var ent T
	return etags.FieldNamesByTags(&ent, etags.TagNamePaginationField, etags.ColumnTagName())
}
