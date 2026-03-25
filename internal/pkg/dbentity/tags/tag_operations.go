// Entity tag operations
package etags

import (
	"reflect"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
)

// TagName is a type alias for a string, used to represent a tag name
type TagName string

const (
	// Default tag name used to specify the database field name.
	// Example: `db:"sale_id"` or `sql:"sale_id"`.
	columnTagNameDefault TagName = "db"

	// Tag name that defines an SQL expression used during SQL query generation.
	// Used for dynamically building SQL expressions.
	// Example: `db:"name" json:"name" sqlexpr:"LOWER(name)"`,
	// which results in SQL: `LOWER(name) = ?`.
	TagNameExpr TagName = "sqlexpr"

	// Tag name controlling inclusion/exclusion of a field in SELECT queries.
	// Allowed values: "false", "true".
	// Default is "true" if not specified.
	// Example: `db:"created_at" select:"false"` — field `created_at` will not be included in SELECT queries.
	TagNameSelect TagName = "select"

	// Tag name controlling inclusion/exclusion of a field in INSERT queries.
	// Allowed values: "false", "true".
	// Default is "true" if not specified.
	// Example: `db:"created_at" insert:"false"` — field `created_at` will not be included in INSERT queries.
	TagNameInsert TagName = "insert"

	// tagValueTrue — tag value indicating that the field is allowed in a specific SQL operation.
	// Used with tags such as `select`, `insert`, etc.
	tagValueTrue = "true"

	// Tag value indicating that the field is forbidden in a specific SQL operation.
	// Used with tags such as `select`, `insert`, etc.
	tagValueFalse = "false"

	// Tag used to mark searchable fields in BaseEntity for Paginate.
	// Allowed values: "true"
	TagNameIsSearch TagName = "issearch"

	// TagNamePaginationField — tag name for specifying a field used in pagination.
	// For example, for sorting or filtering.
	TagNamePaginationField TagName = "pfield"
)

var (
	columnConfig = struct {
		// Tag name used to specify the database field name.
		// Example: `db:"sale_id"` or `sql:"sale_id"`.
		tagName TagName
		once    sync.Once
	}{
		tagName: columnTagNameDefault,
	}
)

// String converts the TagName value into a string
func (t TagName) String() string {
	return string(t)
}

// ColumnTagName returns the current column tag name
func ColumnTagName() TagName {
	return columnConfig.tagName
}

// SetColumnTagName sets the column tag name only once.
// Subsequent calls have no effect.
// Ensures that the tag name is initialized only once during the application lifecycle.
func SetColumnTagName(tagName string) {
	columnConfig.once.Do(func() {
		columnConfig.tagName = TagName(tagName)
	})
}

// Fields returns fields with a tag and without operTagName set to "false".
func Fields(ent any, operTagName TagName) []string {
	return FieldsExceptFalseTag(ent, operTagName)
}

// FieldsWithTrueTag returns a list of fields that have the operTagName explicitly set to tagValueTrue
func FieldsWithTrueTag(ent any, operTagName TagName) []string {
	var fields = make([]string, 0, len(structs.Fields(ent)))

	// Iterate through all struct fields
	for _, field := range structs.Fields(ent) {
		// Get the field name specified in the tag according to columnConfig.tagName
		fieldName := field.Tag(columnConfig.tagName.String())

		// Get the value of the operTagName tag
		operTag := field.Tag(operTagName.String())

		// Condition: field name is specified, and operTagName is set to tagValueTrue
		if fieldName != "" && operTag == tagValueTrue {
			fields = append(fields, fieldName)
		}
	}

	return fields
}

// FieldsExceptFalseTag returns fields with a tag and without operTagName set to "false".
func FieldsExceptFalseTag(ent any, operTagName TagName) []string {
	var fields = make([]string, 0, len(structs.Fields(ent)))

	// Iterate through all struct fields
	for _, field := range structs.Fields(ent) {
		// Get the field name from the corresponding tag
		fieldName := field.Tag(columnConfig.tagName.String())

		if fieldName != "" {
			// Check if the exclusion flag is set in the operTagName tag
			if field.Tag(operTagName.String()) == tagValueFalse {
				continue
			}
			fields = append(fields, fieldName)
		}
	}

	return fields
}

// Values returns field values for which operTagName is not explicitly set to "false".
// If TagNameExpr is set, the value is wrapped with squirrel.Expr.
func Values(ent any, operTagName TagName) []any {
	return ValuesExceptFalseTag(ent, operTagName)
}

// ValuesWithTrueTag returns a list of field values for which the operTagName tag is explicitly set to tagValueTrue.
// If a field has an expression tag (TagNameExpr), its value is wrapped using squirrel.Expr.
func ValuesWithTrueTag(ent any, operTagName TagName) []any {
	var values = make([]any, 0, len(structs.Fields(ent)))

	// Iterate through all struct fields
	for _, field := range structs.Fields(ent) {
		// Get the field name from the corresponding tag
		fieldName := field.Tag(columnConfig.tagName.String())

		if fieldName != "" {
			// Skip if operTagName is not set to tagValueTrue
			if field.Tag(operTagName.String()) != tagValueTrue {
				continue
			}

			// Check if the field has an expression tag (TagNameExpr)
			exprTag := field.Tag(TagNameExpr.String())
			if exprTag != "" {
				// Add the expression (e.g., SQL expression) to the values slice
				values = append(values, squirrel.Expr(exprTag))
			} else {
				// If there is no expression tag, add the field value
				values = append(values, field.Value())
			}
		}
	}

	return values
}

// ValuesExceptFalseTag returns field values for which operTagName is not explicitly set to "false".
// If TagNameExpr is set, the value is wrapped with squirrel.Expr.
func ValuesExceptFalseTag(ent any, operTagName TagName) []any {
	var values = make([]any, 0, len(structs.Fields(ent)))

	// Iterate through all struct fields
	for _, field := range structs.Fields(ent) {
		// Get the field name from the corresponding tag
		fieldName := field.Tag(columnConfig.tagName.String())

		if fieldName != "" {
			// Skip if operTagName is set to tagValueFalse
			if field.Tag(operTagName.String()) == tagValueFalse {
				continue
			}

			// Check if the field has an expression tag (TagNameExpr)
			exprTag := field.Tag(TagNameExpr.String())
			if exprTag != "" {
				// Add the expression (e.g., SQL expression) to the values slice
				values = append(values, squirrel.Expr(exprTag))
			} else {
				// If there is no expression tag, add the field value
				values = append(values, field.Value())
			}
		}
	}

	return values
}

// GetFieldTagByOffset returns the struct field tag for a given field pointer,
// using the field's offset relative to the beginning of the struct.
// Returns the tag value as a string and a boolean indicating whether it was found.
func GetFieldTagByOffset(structPtr any, fieldPtr any, tagKey string) (string, bool) {
	// Get the type of the struct
	structType := reflect.TypeOf(structPtr).Elem()

	// Get the base address of the struct (in bytes)
	baseAddr := reflect.ValueOf(structPtr).Pointer()
	// Get the address of the field (in bytes)
	fieldAddr := reflect.ValueOf(fieldPtr).Pointer()

	// Calculate the field's offset from the beginning of the struct
	offsetAddr := uintptr(fieldAddr - baseAddr)

	// Iterate through all struct fields
	for i := range structType.NumField() {
		field := structType.Field(i)

		// Compare the offset
		if field.Offset == offsetAddr {
			return field.Tag.Get(tagKey), true
		}
	}

	// Field not found
	return "", false
}
