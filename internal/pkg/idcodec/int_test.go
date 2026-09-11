package idcodec

import "testing"

func TestEncodeInt64Base64URL(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  string
	}{
		{
			name:  "zero",
			value: 0,
			want:  "AAAAAAAAAAA",
		},
		{
			name:  "one",
			value: 1,
			want:  "AAAAAAAAAAE",
		},
		{
			name:  "negative",
			value: -1,
			want:  "__________8",
		},
		{
			name:  "timestamp",
			value: 1788706800,
			want:  "AAAAAGqdf_A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeInt64Base64URL(tt.value)

			if got != tt.want {
				t.Errorf("EncodeInt64Base64URL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeInt64Base64URL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int64
	}{
		{
			name:  "zero",
			value: "AAAAAAAAAAA",
			want:  0,
		},
		{
			name:  "one",
			value: "AAAAAAAAAAE",
			want:  1,
		},
		{
			name:  "negative",
			value: "__________8",
			want:  -1,
		},
		{
			name:  "timestamp",
			value: "AAAAAGqdf_A",
			want:  1788706800,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeInt64Base64URL(tt.value)
			if err != nil {
				t.Fatalf("DecodeInt64Base64URL() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("DecodeInt64Base64URL() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEncodeDecodeInt64Base64URL(t *testing.T) {
	values := []int64{
		0,
		1,
		-1,
		1788706800,
		-1788706800,
	}

	for _, value := range values {
		t.Run("round trip", func(t *testing.T) {
			encoded := EncodeInt64Base64URL(value)

			got, err := DecodeInt64Base64URL(encoded)
			if err != nil {
				t.Fatalf("DecodeInt64Base64URL() error = %v", err)
			}

			if got != value {
				t.Errorf("round trip = %d, want %d", got, value)
			}
		})
	}
}
