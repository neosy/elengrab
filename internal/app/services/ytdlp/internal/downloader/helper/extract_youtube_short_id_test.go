package helper

import "testing"

func TestExtractYouTubeShortID(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantID  string
		wantErr bool
	}{
		{
			name:    "valid shorts url",
			rawURL:  "https://www.youtube.com/shorts/Ulpa7scJnPw",
			wantID:  "Ulpa7scJnPw",
			wantErr: false,
		},
		{
			name:    "valid shorts url with query",
			rawURL:  "https://www.youtube.com/shorts/Ulpa7scJnPw?si=abc123",
			wantID:  "Ulpa7scJnPw",
			wantErr: false,
		},
		{
			name:    "valid shorts url mobile",
			rawURL:  "https://m.youtube.com/shorts/Ulpa7scJnPw",
			wantID:  "Ulpa7scJnPw",
			wantErr: false,
		},
		{
			name:    "not shorts url",
			rawURL:  "https://www.youtube.com/watch?v=Ulpa7scJnPw",
			wantErr: true,
		},
		{
			name:    "missing short id",
			rawURL:  "https://www.youtube.com/shorts/",
			wantErr: true,
		},
		{
			name:    "invalid short id characters",
			rawURL:  "https://www.youtube.com/shorts/Ulpa7scJnP!",
			wantErr: true,
		},
		{
			name:    "invalid short id length",
			rawURL:  "https://www.youtube.com/shorts/abc",
			wantErr: true,
		},
		{
			name:    "extra path segments",
			rawURL:  "https://www.youtube.com/shorts/Ulpa7scJnPw/extra",
			wantErr: true,
		},
		{
			name:    "invalid host",
			rawURL:  "https://google.com/shorts/Ulpa7scJnPw",
			wantErr: true,
		},
		{
			name:    "not a url",
			rawURL:  "%%%bad-url%%%",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := ExtractYouTubeShortID(tt.rawURL)

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
