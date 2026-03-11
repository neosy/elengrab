package handlers

var (
	mapContentTypeByExt = map[string]string{
		"mp3":  "audio/mpeg",
		"m4a":  "audio/mp4",
		"aac":  "audio/aac",
		"ogg":  "audio/ogg",
		"opus": "audio/opus",
		"flac": "audio/flac",

		"mp4":  "video/mp4",
		"webm": "video/webm",
		"mkv":  "video/x-matroska",
		"mov":  "video/quicktime",
	}
)
