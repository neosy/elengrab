package stringx

import (
	"testing"
)

func TestTruncateBytesWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxBytes int
		want     string
	}{
		{
			name:     "does not truncate when string fits",
			input:    "Hello world",
			maxBytes: 20,
			want:     "Hello world",
		},
		{
			name:     "returns empty for empty string",
			input:    "",
			maxBytes: 10,
			want:     "",
		},
		{
			name:     "keeps fully included last word",
			input:    "Hello world test",
			maxBytes: 11,
			want:     "Hello world",
		},
		{
			name:     "removes partially truncated last word",
			input:    "Hello world test",
			maxBytes: 12,
			want:     "Hello world",
		},
		{
			name:     "handles cyrillic",
			input:    "Привет мир тест",
			maxBytes: 19,
			want:     "Привет мир",
		},
		{
			name:     "does not break utf8 character",
			input:    "Привет мир",
			maxBytes: 10,
			want:     "Приве",
		},
		{
			name:     "returns first word when it fully fits",
			input:    "Привет мир",
			maxBytes: 13,
			want:     "Привет",
		},
		{
			name:     "truncates before word when next word does not fit",
			input:    "Hello world",
			maxBytes: 7,
			want:     "Hello",
		},
		{
			name:     "handles first word longer than max bytes",
			input:    "Superlongword test",
			maxBytes: 5,
			want:     "Super",
		},
		{
			name:     "keeps word when limit ends exactly at word boundary",
			input:    "Hello world",
			maxBytes: 5,
			want:     "Hello",
		},
		{
			name:     "handles multiple spaces",
			input:    "Hello   world",
			maxBytes: 8,
			want:     "Hello",
		},
		{
			name:     "handles tabs as spaces",
			input:    "Hello\tworld",
			maxBytes: 6,
			want:     "Hello",
		},
		{
			name:     "returns empty when maxBytes is zero",
			input:    "Hello world",
			maxBytes: 0,
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateBytesWords(tt.input, tt.maxBytes)

			if got != tt.want {
				t.Errorf("TruncateBytesWords() = %q, want %q", got, tt.want)
			}
		})
	}
}
