package database

import (
	"runtime"
	"testing"
)

func TestFilePathForMigrateWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "absolute path",
			path: `C:\migrations`,
			want: "file:///C:/migrations",
		},
		{
			name: "nested path",
			path: `C:\Users\user\AppData\Local\elengrab\migrations`,
			want: "file:///C:/Users/user/AppData/Local/elengrab/migrations",
		},
		{
			name: "path with spaces",
			path: `C:\Program Files\elengrab\migrations`,
			want: "file:///C:/Program%20Files/elengrab/migrations",
		},
		{
			name: "path with special characters",
			path: `C:\temp\my migrations\test#1`,
			want: "file:///C:/temp/my%20migrations/test%231",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filePathForMigrate(tt.path)

			if got != tt.want {
				t.Errorf("filePathForMigrate(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestFilePathForMigrateUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "absolute path",
			path: "/tmp/migrations",
			want: "file:///tmp/migrations",
		},
		{
			name: "nested path",
			path: "/var/lib/elengrab/migrations",
			want: "file:///var/lib/elengrab/migrations",
		},
		{
			name: "path with spaces",
			path: "/tmp/my migrations/test",
			want: "file:///tmp/my migrations/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filePathForMigrate(tt.path)

			if got != tt.want {
				t.Errorf("filePathForMigrate(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
