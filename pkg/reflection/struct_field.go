package reflection

import (
	"reflect"
)

// StructFieldName returns the field name from the tag by the struct field pointer
// Example:
// var tradeEmpty = tradeEntity{}
// StructFieldName(&tradeEmpty, &tradeEmpty.TradeId, "sql")
func StructFieldName(structPtr any, fieldPtr any, tag string) (string, error) {
	// Get the value via reflect
	structVal := reflect.ValueOf(structPtr)

	// Check the type of the first argument (pointer to a struct)
	if structVal.Kind() != reflect.Pointer || structVal.Elem().Kind() != reflect.Struct {
		return "", ErrFirstArgumentTypeMustPointerStructure
	}

	// Check the type of the second argument (pointer to a field)
	fieldVal := reflect.ValueOf(fieldPtr)
	if fieldVal.Kind() != reflect.Ptr {
		return "", ErrSecondArgumentTypeMustPointerFieldStructure
	}

	var name string

	// Get the type of the struct via reflect
	structType := reflect.ValueOf(structPtr).Elem().Type()

	// Iterate through all struct fields
	for i := range structType.NumField() {
		// Get the struct field
		structField := structType.Field(i)

		// If this is an anonymous struct, check its fields recursively
		if structField.Anonymous {
			// Get the value of the anonymous struct
			anonymousFieldValue := reflect.ValueOf(structPtr).Elem().Field(i)
			// Check: if it is a struct, search the field recursively
			if anonymousFieldValue.Kind() == reflect.Struct {
				// Recursively search the field in the anonymous struct
				subName, err := StructFieldName(anonymousFieldValue.Addr().Interface(), fieldPtr, tag)
				if err == nil && subName != "" {
					name = subName
					break
				} else {
					continue
				}
			} else {
				continue
			}
		}

		// Check if the pointer to the field matches the address of the current struct field
		if reflect.ValueOf(fieldPtr).Pointer() == reflect.ValueOf(structPtr).Elem().Field(i).Addr().Pointer() {
			// Return the value from the tag (e.g., "sql")
			name = structType.Field(i).Tag.Get(tag)
			break
		}
	}

	if name == "" {
		return "", ErrStructureFieldNameEmpty
	}

	return name, nil
}
