package reflection

import (
	"reflect"
	"testing"
)

func TestStructToMap_AllCombinations(t *testing.T) {

	type Mixed struct {
		Int     int     `json:"int_field"`
		Bool    bool    `json:"bool_field"`
		Float   float64 `json:"float_field"`
		String  string  `json:"string_field"`
		Pointer *string `json:"ptr_field"`
		Skip    string  `json:"-"`
		NoTag   int
	}

	t.Run("full struct with all types", func(t *testing.T) {
		str := "hello"

		in := Mixed{
			Int:     42,
			Bool:    true,
			Float:   3.14,
			String:  "text",
			Pointer: &str,
			Skip:    "hidden",
			NoTag:   7,
		}

		got := StructToMap(in)

		want := map[string]any{
			"int_field":    42,
			"bool_field":   true,
			"float_field":  3.14,
			"string_field": "text",
			"ptr_field":    &str,
			"NoTag":        7,
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch\nwant: %#v\ngot:  %#v", want, got)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		str := "p"

		in := &Mixed{
			Int:     1,
			Bool:    false,
			Float:   1.1,
			String:  "p",
			Pointer: &str,
			NoTag:   99,
		}

		got := StructToMap(in)

		want := map[string]any{
			"int_field":    1,
			"bool_field":   false,
			"float_field":  1.1,
			"string_field": "p",
			"ptr_field":    &str,
			"NoTag":        99,
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch\nwant: %#v\ngot:  %#v", want, got)
		}
	})

	t.Run("zero values included", func(t *testing.T) {
		in := Mixed{}

		got := StructToMap(in)

		want := map[string]any{
			"int_field":    0,
			"bool_field":   false,
			"float_field":  0.0,
			"string_field": "",
			"ptr_field":    (*string)(nil),
			"NoTag":        0,
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch\nwant: %#v\ngot:  %#v", want, got)
		}
	})

	t.Run("non struct input", func(t *testing.T) {
		got := StructToMap(123)

		if len(got) != 0 {
			t.Fatalf("expected empty map, got %#v", got)
		}
	})

	t.Run("json tag override priority", func(t *testing.T) {
		type S struct {
			Field string `json:"custom_name"`
		}

		in := S{Field: "value"}

		got := StructToMap(in)

		want := map[string]any{
			"custom_name": "value",
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("mismatch\nwant: %#v\ngot:  %#v", want, got)
		}
	})
}
