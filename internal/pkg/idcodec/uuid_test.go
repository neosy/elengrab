package idcodec

import (
	"testing"

	"github.com/google/uuid"
)

func TestEncodeUUIDBase64URL(t *testing.T) {
	tests := []struct {
		name string
		id   uuid.UUID
		want string
	}{
		{
			name: "zero",
			id:   uuid.Nil,
			want: "AAAAAAAAAAAAAAAAAAAAAA",
		},
		{
			name: "uuid",
			id:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			want: "VQ6EAOKbQdSnFkRmVUQAAA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeUUIDBase64URL(tt.id)

			if got != tt.want {
				t.Errorf("EncodeUUIDBase64URL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeUUIDBase64URL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  uuid.UUID
	}{
		{
			name:  "zero",
			value: "AAAAAAAAAAAAAAAAAAAAAA",
			want:  uuid.Nil,
		},
		{
			name:  "uuid",
			value: "VQ6EAOKbQdSnFkRmVUQAAA",
			want:  uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeUUIDBase64URL(tt.value)
			if err != nil {
				t.Fatalf("DecodeUUIDBase64URL() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("DecodeUUIDBase64URL() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDecodeUUIDBase64URLInvalid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "invalid base64",
			value: "!!!",
		},
		{
			name:  "invalid uuid length",
			value: "AQ",
		},
		{
			name:  "too long",
			value: "AAAAAAAAAAAAAAAAAAAAAAA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeUUIDBase64URL(tt.value)
			if err == nil {
				t.Fatal("DecodeUUIDBase64URL() error = nil, want error")
			}
		})
	}
}

func TestEncodeDecodeUUIDBase64URL(t *testing.T) {
	ids := []uuid.UUID{
		uuid.Nil,
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
	}

	for _, id := range ids {
		t.Run(id.String(), func(t *testing.T) {
			encoded := EncodeUUIDBase64URL(id)

			got, err := DecodeUUIDBase64URL(encoded)
			if err != nil {
				t.Fatalf("DecodeUUIDBase64URL() error = %v", err)
			}

			if got != id {
				t.Errorf("DecodeUUIDBase64URL() = %s, want %s", got, id)
			}
		})
	}
}
