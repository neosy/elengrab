package reflection

import "testing"

func BenchmarkStructToMap_AllCases(b *testing.B) {

	b.Run("mixed types struct", func(b *testing.B) {

		type Mixed struct {
			Int     int     `json:"int_field"`
			Bool    bool    `json:"bool_field"`
			Float   float64 `json:"float_field"`
			String  string  `json:"string_field"`
			Pointer *string `json:"ptr_field"`
			NoTag   int
			Skip    string `json:"-"`
		}

		str := "hello"
		in := Mixed{
			Int:     10,
			Bool:    true,
			Float:   3.14,
			String:  "text",
			Pointer: &str,
			NoTag:   99,
			Skip:    "hidden",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = StructToMap(in)
		}
	})

	b.Run("pointer struct", func(b *testing.B) {

		type S struct {
			A int
			B string `json:"b_field"`
			C bool
		}

		in := &S{
			A: 1,
			B: "value",
			C: true,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = StructToMap(in)
		}
	})

	b.Run("zero values heavy struct", func(b *testing.B) {

		type Z struct {
			A int
			B int
			C int
			D int
			E int
			F string
			G bool
			H float64
		}

		in := Z{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = StructToMap(in)
		}
	})

	b.Run("many fields struct", func(b *testing.B) {

		type Big struct {
			F1  int
			F2  int
			F3  int
			F4  int
			F5  int
			F6  int
			F7  int
			F8  int
			F9  int
			F10 int
			F11 int
			F12 int
			F13 int
			F14 int
			F15 int
		}

		in := Big{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = StructToMap(in)
		}
	})
}
