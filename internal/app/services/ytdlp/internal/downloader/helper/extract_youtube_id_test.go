package helper

import "testing"

func TestExtractYouTubeID(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantID  string
		wantErr bool
	}{
		{
			name:    "youtu.be url",
			rawURL:  "https://youtu.be/Ulpa7scJnPw",
			wantID:  "Ulpa7scJnPw",
			wantErr: false,
		},
		{
			name:    "youtu.be url with query",
			rawURL:  "https://youtu.be/Ulpa7scJnPw?si=abc123",
			wantID:  "Ulpa7scJnPw",
			wantErr: false,
		},
		{
			name:    "youtube watch url",
			rawURL:  "https://www.youtube.com/watch?v=6WGwkQfBTgM",
			wantID:  "6WGwkQfBTgM",
			wantErr: false,
		},
		{
			name:    "youtube watch url with extra params",
			rawURL:  "https://www.youtube.com/watch?v=6WGwkQfBTgM&t=42",
			wantID:  "6WGwkQfBTgM",
			wantErr: false,
		},
		{
			name:    "invalid host",
			rawURL:  "https://google.com/watch?v=6WGwkQfBTgM",
			wantErr: true,
		},
		{
			name:    "missing video id",
			rawURL:  "https://www.youtube.com/watch",
			wantErr: true,
		},
		{
			name:    "invalid id length",
			rawURL:  "https://youtu.be/abc",
			wantErr: true,
		},
		{
			name:    "invalid id character",
			rawURL:  "https://youtu.be/Ulpa7scJnP!",
			wantErr: true,
		},
		{
			name:    "invalid path",
			rawURL:  "https://youtu.be/foo/bar",
			wantErr: true,
		},
		{
			name:    "not url",
			rawURL:  "%%%bad-url%%%",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := ExtractYouTubeID(tt.rawURL)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotID != tt.wantID {
				t.Fatalf("expected id %q, got %q", tt.wantID, gotID)
			}
		})
	}
}
