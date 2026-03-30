package reflection

import (
	"reflect"
)

// StructFieldPointers returns a slice of pointers to all **exported** fields of a struct
func StructFieldPointers(structPtr any) ([]any, error) {
	val := reflect.ValueOf(structPtr)

	// Check that input is a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, ErrInputMustBePointerToStruct
	}

	val = val.Elem()
	numFields := val.NumField()
	ptrs := make([]any, 0, numFields)

	// Iterate over all fields
	for i := range numFields {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		// Skip anonymous (embedded) fields
		if fieldType.Anonymous {
			continue
		}

		// Include only exported fields (fields with uppercase first letter)
		if field.CanAddr() {
			ptrs = append(ptrs, field.Addr().Interface())
		}
	}

	return ptrs, nil
}
