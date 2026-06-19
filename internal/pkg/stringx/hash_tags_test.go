package stringx

import (
	"testing"
)

func TestRemoveHashtags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes hashtags from text",
			input: "Funny parrot video #ParrotTeacher #CutePets",
			want:  "Funny parrot video",
		},
		{
			name:  "removes multiple hashtags",
			input: "#Tag1 #Tag2 #Tag3",
			want:  "",
		},
		{
			name:  "keeps text without hashtags",
			input: "Funny parrot video",
			want:  "Funny parrot video",
		},
		{
			name:  "removes hashtags in the middle",
			input: "Funny #Parrot video #CutePets today",
			want:  "Funny video today",
		},
		{
			name:  "cleans extra spaces after removing",
			input: "Funny   video   #Tag",
			want:  "Funny video",
		},
		{
			name:  "supports unicode hashtags",
			input: "Видео с попугаем #Попугай",
			want:  "Видео с попугаем",
		},
		{
			name:  "keeps empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveHashtags(tt.input)

			if got != tt.want {
				t.Errorf("RemoveHashtags() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoveTrailingHashtags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "removes trailing hashtags",
			input: "Funny parrot video #ParrotTeacher #CutePets",
			want:  "Funny parrot video",
		},
		{
			name:  "removes hashtag without space",
			input: "Coin Flicking Challenge#funny",
			want:  "Coin Flicking Challenge",
		},
		{
			name:  "removes mixed hashtags",
			input: "Coin Flicking Challenge#funny #CutePets#funny",
			want:  "Coin Flicking Challenge",
		},
		{
			name:  "removes mixed hashtags",
			input: "Funny parrot video#ParrotTeacher #CutePets #ViralReel",
			want:  "Funny parrot video",
		},
		{
			name:  "keeps hashtags in middle",
			input: "Funny #Parrot video",
			want:  "Funny #Parrot video",
		},
		{
			name:  "keeps text without hashtags",
			input: "Funny parrot video",
			want:  "Funny parrot video",
		},
		{
			name:  "handles only hashtag",
			input: "#ParrotTeacher",
			want:  "",
		},
		{
			name:  "handles empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveTrailingHashtags(tt.input)

			if got != tt.want {
				t.Errorf("RemoveTrailingHashtags() = %q, want %q", got, tt.want)
			}
		})
	}
}
