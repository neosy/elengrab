package reflection

import (
	"reflect"
)

// StructFieldPointer returns a pointer to a struct's field specified by name or tag.
// If tag is empty, it searches by field name. Otherwise, it searches for a field with the given tag key.
func StructFieldPointer(structPtr any, fieldName string, tag string) (any, error) {
	val := reflect.ValueOf(structPtr)

	// Check that input is a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, ErrInputMustBePointerToStruct
	}

	val = val.Elem()
	typ := val.Type()
	numFields := val.NumField()

	for i := range numFields {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip anonymous fields
		if fieldType.Anonymous {
			continue
		}

		// Only exported fields
		if !field.CanAddr() {
			continue
		}

		// If tag is provided, match by tag value
		if tag != "" {
			tagValue := fieldType.Tag.Get(tag)
			if tagValue != "" && tagValue == fieldName {
				return field.Addr().Interface(), nil
			}
		} else {
			// Match by field name
			if fieldType.Name == fieldName {
				return field.Addr().Interface(), nil
			}
		}
	}

	return nil, ErrFieldNotFound
}
