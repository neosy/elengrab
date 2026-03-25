package dbentity

import (
	etags "github.com/neosy/elengrab/internal/pkg/dbentity/tags"
	"github.com/neosy/elengrab/internal/pkg/reflection"
)

// BaseEntity is a base structure that provides common methods
// for working with entity fields and their SQL tag names
type BaseEntity[T any] struct{}

// FieldName field name from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId)
func (e *BaseEntity[T]) FieldName(structPtr *T, fieldPtr any) string {
	return e.FieldNameWithAlias(structPtr, fieldPtr, "")
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *BaseEntity[T]) FieldNameWithAlias(structPtr *T, fieldPtr any, alias string) string {
	name, _ := reflection.StructFieldName(structPtr, fieldPtr, etags.ColumnTagName().String())

	if alias == "" {
		return name
	}

	return alias + "." + name
}

// PaginateFieldName returns the field name with alias from the pagination tag (`pfield`)
// based on the pointer to the struct and the pointer to the field.
// Example:
// var ent <TableEntity>
// ent.PaginateFieldName(&ent, &ent.SalesId)
func (e *BaseEntity[T]) PaginateFieldName(structPtr *T, fieldPtr any) string {
	name, _ := reflection.StructFieldName(structPtr, fieldPtr, etags.TagNamePaginationField.String())
	return name
}

// FieldsAll returns a list of fields that will be used for queries
func (e *BaseEntity[T]) FieldsAll() []string {
	var ent T
	return etags.Fields(&ent, etags.TagNameSelect)
}

// SearchableFields returns a list of fields explicitly marked as searchable
func (e *BaseEntity[T]) SearchableFields() []string {
	var ent T
	return etags.FieldsWithTrueTag(&ent, etags.TagNameIsSearch)
}

// FieldsAllWithAlias returns a list of fields with alias that will be used for queries
func (e *BaseEntity[T]) FieldsAllWithAlias(alias string) []string {
	fields := e.FieldsAll()

	withAlias := make([]string, len(fields))
	for i, f := range fields {
		withAlias[i] = alias + "." + f
	}

	return withAlias
}

// Fields returns a list of fields allowed for insert or update operations
// where the tag 'insert' tag is not explicitly set to "false".
func (e *BaseEntity[T]) Fields() []string {
	var ent T
	return etags.FieldsExceptFalseTag(&ent, etags.TagNameInsert)
}

// Values returns field values allowed for insert or update operations
// where the tag 'insert' tag is not explicitly set to "false".
func (e *BaseEntity[T]) Values(structPtr *T) []any {
	return etags.ValuesExceptFalseTag(structPtr, etags.TagNameInsert)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
// This can be used for scanning database rows into the struct.
func (e *BaseEntity[T]) FieldPointers(structPtr *T) ([]any, error) {
	return reflection.StructFieldPointers(structPtr)
}
