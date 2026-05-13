package reflection

import (
	"reflect"
	"strings"
)

// StructToMap converts a struct to map[string]any.
// It uses the json tag as the key if present.
// If json tag is missing, field name is used.
// Fields with json:"-" are skipped.
//
// Performance (AMD Ryzen 5 5600G):
// ----------------------------------
// mixed struct        : ~888.3 ns/op   | ~1.13M ops/sec | 512 B/op  | 9 allocs/op
// pointer struct      : ~475.6 ns/op   | ~2.10M ops/sec | 416 B/op  | 8 allocs/op
// zero values struct  : ~693.3 ns/op   | ~1.44M ops/sec | 464 B/op  | 10 allocs/op
// many fields (15)    : ~2570 ns/op    | ~0.39M ops/sec | 2376 B/op | 22 allocs/op
//
// Notes:
// - Performance is dominated by reflect.Value.Interface() allocations.
// - Complexity is O(n) with high constant overhead per field.
// - Suitable for low/medium-frequency usage (logging, DTO mapping, debug layers).
// - Not recommended for high-QPS or hot-path execution.
func StructToMap(v any) map[string]any {
	result := make(map[string]any)

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Dereference pointer if needed
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Return empty map if value is not a struct
	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Read json tag
		tag := field.Tag.Get("json")

		// Skip ignored fields
		if tag == "-" {
			continue
		}

		// Remove options like omitempty
		tagParts := strings.Split(tag, ",")

		var key string

		// Use json tag if exists
		if tagParts[0] != "" {
			key = tagParts[0]
		} else {
			key = field.Name
		}

		result[key] = value.Interface()
	}

	return result
}
