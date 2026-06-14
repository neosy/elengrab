package stringx

import "testing"

func TestLowerFirstWord(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "multiple words",
			input:    "ACCESS DENIED Error",
			expected: "access DENIED Error",
		},
		{
			name:     "already lower first word",
			input:    "access DENIED Error",
			expected: "access DENIED Error",
		},
		{
			name:     "single word with mixed case",
			input:    "TeSt",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LowerFirstWord(tt.input)
			if got != tt.expected {
				t.Errorf("LowerFirstWord(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
