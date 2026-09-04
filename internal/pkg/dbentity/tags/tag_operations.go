// Entity tag operations
package etags

import (
	"reflect"
	"sync"

	"github.com/Masterminds/squirrel"
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

// FieldNames returns fields with a tag and without operationTagName set to "false".
func FieldNames(ent any, operationTagName TagName) []string {
	return FieldNamesExceptFalseTag(ent, operationTagName)
}

// FieldsWithTrueTag returns a list of fields that have the operTagName explicitly set to tagValueTrue.
func FieldsWithTrueTag(ent any, operTagName TagName) []string {
	typeInfo := structType(ent)

	result := make([]string, 0, typeInfo.NumField())

	for field := range typeInfo.Fields() {
		// Get the field name specified in the tag according to columnConfig.tagName
		fieldName := field.Tag.Get(columnConfig.tagName.String())

		// Get the value of the operTagName tag
		operTag := field.Tag.Get(operTagName.String())

		// Include fields explicitly enabled for the operation.
		if fieldName != "" && operTag == tagValueTrue {
			result = append(result, fieldName)
		}
	}

	return result
}

// FieldNamesExceptFalseTag returns fields with a tag and without operationTagName set to "false".
func FieldNamesExceptFalseTag(ent any, operationTagName TagName) []string {
	typeInfo := structType(ent)

	result := make([]string, 0, typeInfo.NumField())

	for field := range typeInfo.Fields() {
		fieldName := field.Tag.Get(columnConfig.tagName.String())

		if fieldName != "" &&
			field.Tag.Get(operationTagName.String()) != tagValueFalse {
			result = append(result, fieldName)
		}
	}

	return result
}

// FieldNamesByTag returns a list of field names defined by the specified tag.
// Only fields with a non-empty value for the specified tag are included.
func FieldNamesByTag(ent any, tag TagName) []string {
	typeInfo := structType(ent)

	names := make([]string, 0)

	for field := range typeInfo.Fields() {
		fieldName := field.Tag.Get(tag.String())

		if fieldName != "" {
			names = append(names, fieldName)
		}
	}

	return names
}

// FieldNamesByTags returns a map that maps values of tagA to values of tagB.
// Only fields that have both tags defined are included.
//
// For example, given:
//
//	SourceCreatedAt time.Time `db:"created_at" pfield:"createdAt"`
//
// FieldNamesByTags(ent, "pfield", "db") returns:
//
//	map[string]string{"createdAt": "created_at"}
func FieldNamesByTags(ent any, tagA, tagB TagName) map[string]string {
	typeInfo := structType(ent)

	names := make(map[string]string)

	for field := range typeInfo.Fields() {
		fieldAName := field.Tag.Get(tagA.String())
		fieldBName := field.Tag.Get(tagB.String())

		if fieldAName != "" && fieldBName != "" {
			names[fieldAName] = fieldBName
		}
	}

	return names
}

// Values returns field values for which operTagName is not explicitly set to "false".
// If TagNameExpr is set, the value is wrapped with squirrel.Expr.
func Values(ent any, operTagName TagName) []any {
	return ValuesExceptFalseTag(ent, operTagName)
}

// ValuesWithTrueTag returns a list of field values for which the operTagName tag is explicitly set to tagValueTrue.
// If a field has an expression tag (TagNameExpr), its value is wrapped using squirrel.Expr.
func ValuesWithTrueTag(ent any, operTagName TagName) []any {
	typeInfo := structType(ent)
	value := structValue(ent)

	values := make([]any, 0, typeInfo.NumField())

	for i := range typeInfo.NumField() {
		field := typeInfo.Field(i)

		// Skip fields without a column tag.
		if field.Tag.Get(columnConfig.tagName.String()) == "" {
			continue
		}

		// Skip if operTagName is not set to tagValueTrue.
		if field.Tag.Get(operTagName.String()) != tagValueTrue {
			continue
		}

		// Use SQL expression when specified.
		if exprTag := field.Tag.Get(TagNameExpr.String()); exprTag != "" {
			values = append(values, squirrel.Expr(exprTag))
			continue
		}

		values = append(values, value.Field(i).Interface())
	}

	return values
}

// ValuesExceptFalseTag returns field values for which operTagName is not explicitly set to "false".
// If TagNameExpr is set, the value is wrapped with squirrel.Expr.
func ValuesExceptFalseTag(ent any, operTagName TagName) []any {
	typeInfo := structType(ent)
	value := structValue(ent)

	values := make([]any, 0, typeInfo.NumField())

	for i := range typeInfo.NumField() {
		field := typeInfo.Field(i)

		// Skip fields without a column tag.
		if field.Tag.Get(columnConfig.tagName.String()) == "" {
			continue
		}

		// Skip fields explicitly excluded from the operation.
		if field.Tag.Get(operTagName.String()) == tagValueFalse {
			continue
		}

		// Use SQL expression when specified.
		if exprTag := field.Tag.Get(TagNameExpr.String()); exprTag != "" {
			values = append(values, squirrel.Expr(exprTag))
			continue
		}

		values = append(values, value.Field(i).Interface())
	}

	return values
}

// GetFieldTagByOffset returns the struct field tag for a given field pointer,
// using the field's offset relative to the beginning of the struct.
// Returns the tag value as a string and a boolean indicating whether it was found.
func GetFieldTagByOffset(structPtr any, fieldPtr any, tagKey string) (string, bool) {
	typeInfo := structType(structPtr)

	// Get the base address of the struct (in bytes)
	baseAddr := reflect.ValueOf(structPtr).Pointer()

	// Get the address of the field (in bytes)
	fieldAddr := reflect.ValueOf(fieldPtr).Pointer()

	// Calculate the field's offset from the beginning of the struct
	offsetAddr := uintptr(fieldAddr - baseAddr)

	// Iterate through all struct fields
	for field := range typeInfo.Fields() {
		if field.Offset == offsetAddr {
			return field.Tag.Get(tagKey), true
		}
	}

	return "", false
}

func structType(ent any) reflect.Type {
	typeInfo := reflect.TypeOf(ent)

	if typeInfo.Kind() == reflect.Pointer {
		typeInfo = typeInfo.Elem()
	}

	return typeInfo
}

func structValue(ent any) reflect.Value {
	value := reflect.ValueOf(ent)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	return value
}
