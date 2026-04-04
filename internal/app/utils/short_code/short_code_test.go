package shcode

import (
	"testing"
)

// GenerateShortCode tests the GenerateShortCode function.
func TestGenerateShortCode(t *testing.T) {
	// Test cases
	tests := []struct {
		id             string
		length         uint8
		expectedLength uint8
		expectedErr    bool
	}{
		// Test 1: Empty string (expect empty string as result)
		{"", 8, 0, false},

		// Test 2: String with length 8 (should return a code of length 8)
		{"test-id", 8, 8, false},

		// Test 3: String with length 64 (check full base64 string length)
		{"longer-test-id", 64, 40, false}, // base64 string will be length 40

		// Test 4: Length greater than base64 length (ensure it is not truncated incorrectly)
		{"another-id", 100, 40, false},

		// Test 5: String with UUID, URL, and timestamp
		{
			id:             "550e8400-e29b-41d4-a716-446655440000:https://example.com:1713200000000000000",
			length:         6,
			expectedLength: 6,
			expectedErr:    false,
		},

		// Test 6: String with UUID, URL, and timestamp
		{
			id:             "123e4567-e89b-12d3-a456-426614174000:https://google.com:1713201234567890000",
			length:         12,
			expectedLength: 12,
			expectedErr:    false,
		},

		// Test 7: String with UUID, URL, and timestamp
		{
			id:             "f47ac10b-58cc-4372-a567-0e02b2c3d479:https://openai.com:1713219876543210000",
			length:         16,
			expectedLength: 16,
			expectedErr:    false,
		},

		// Test 8: String with UUID, URL, and timestamp
		{
			id:             "9c858901-8a57-4791-81fe-4c455b099bc9:https://example.org:1713222222222222222",
			length:         0, // no limit
			expectedLength: 40,
			expectedErr:    false,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			// Generate code using the function
			got := GenerateShortCode(tt.id, tt.length)

			// Check for empty string when id is empty
			if tt.expectedErr && got != "" {
				t.Errorf("GenUrlCode() = %v, want empty string", got)
			}

			// Check result length
			if len(got) != int(tt.expectedLength) {
				t.Errorf("GenUrlCode() = %v, want length = %d", got, tt.expectedLength)
			}
		})
	}
}
