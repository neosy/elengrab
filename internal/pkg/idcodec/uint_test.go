package idcodec

import "testing"

func TestEncodeUintBase64URL(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name:  "uint8",
			value: uint8(255),
			want:  "_w",
		},
		{
			name:  "uint16",
			value: uint16(65535),
			want:  "__8",
		},
		{
			name:  "uint32",
			value: uint32(1),
			want:  "AAAAAQ",
		},
		{
			name:  "uint64",
			value: uint64(1),
			want:  "AAAAAAAAAAE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string

			switch value := tt.value.(type) {
			case uint8:
				got = EncodeUintBase64URL(value)
			case uint16:
				got = EncodeUintBase64URL(value)
			case uint32:
				got = EncodeUintBase64URL(value)
			case uint64:
				got = EncodeUintBase64URL(value)
			default:
				t.Fatalf("unsupported type %T", tt.value)
			}

			if got != tt.want {
				t.Errorf("EncodeUintBase64URL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeUintBase64URL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  any
	}{
		{
			name:  "uint8",
			value: "_w",
			want:  uint8(255),
		},
		{
			name:  "uint16",
			value: "__8",
			want:  uint16(65535),
		},
		{
			name:  "uint32",
			value: "AAAAAQ",
			want:  uint32(1),
		},
		{
			name:  "uint64",
			value: "AAAAAAAAAAE",
			want:  uint64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch want := tt.want.(type) {
			case uint8:
				got, err := DecodeUintBase64URL[uint8](tt.value)
				if err != nil {
					t.Fatalf("DecodeUintBase64URL() error = %v", err)
				}

				if got != want {
					t.Errorf("DecodeUintBase64URL() = %d, want %d", got, want)
				}

			case uint16:
				got, err := DecodeUintBase64URL[uint16](tt.value)
				if err != nil {
					t.Fatalf("DecodeUintBase64URL() error = %v", err)
				}

				if got != want {
					t.Errorf("DecodeUintBase64URL() = %d, want %d", got, want)
				}

			case uint32:
				got, err := DecodeUintBase64URL[uint32](tt.value)
				if err != nil {
					t.Fatalf("DecodeUintBase64URL() error = %v", err)
				}

				if got != want {
					t.Errorf("DecodeUintBase64URL() = %d, want %d", got, want)
				}

			case uint64:
				got, err := DecodeUintBase64URL[uint64](tt.value)
				if err != nil {
					t.Fatalf("DecodeUintBase64URL() error = %v", err)
				}

				if got != want {
					t.Errorf("DecodeUintBase64URL() = %d, want %d", got, want)
				}
			}
		})
	}
}

func TestDecodeUintBase64URLInvalid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "invalid base64",
			value: "!!!",
		},
		{
			name:  "invalid uint8 length",
			value: "AAAAAQ",
		},
		{
			name:  "invalid uint16 length",
			value: "AAAAAQ",
		},
		{
			name:  "invalid uint32 length",
			value: "AAAAAAAAAAE",
		},
		{
			name:  "invalid uint64 length",
			value: "AAAAAQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			switch tt.name {
			case "invalid base64":
				_, err = DecodeUintBase64URL[uint32](tt.value)
			case "invalid uint8 length":
				_, err = DecodeUintBase64URL[uint8](tt.value)
			case "invalid uint16 length":
				_, err = DecodeUintBase64URL[uint16](tt.value)
			case "invalid uint32 length":
				_, err = DecodeUintBase64URL[uint32](tt.value)
			case "invalid uint64 length":
				_, err = DecodeUintBase64URL[uint64](tt.value)
			}

			if err == nil {
				t.Fatal("DecodeUintBase64URL() error = nil, want error")
			}
		})
	}
}
