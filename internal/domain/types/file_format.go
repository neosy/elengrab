package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FileFormat string

const (
	FileFormatNone FileFormat = "none"
	FileFormatAuto FileFormat = "auto"

	FileFormatMP4  FileFormat = "mp4"
	FileFormatWebM FileFormat = "webm"

	FileFormatMP3  FileFormat = "mp3"
	FileFormatM4A  FileFormat = "m4a"
	FileFormatFLAC FileFormat = "flac"
	FileFormatOPUS FileFormat = "opus"
)

var (
	fileFormats = map[FileFormat]struct{}{
		FileFormatAuto: {},
		FileFormatMP4:  {},
		FileFormatWebM: {},
		FileFormatMP3:  {},
		FileFormatM4A:  {},
		FileFormatFLAC: {},
		FileFormatOPUS: {},
	}

	extNoDotByFileFormat = map[FileFormat]string{
		FileFormatMP4:  "mp4",
		FileFormatWebM: "webm",
		FileFormatMP3:  "mp3",
		FileFormatM4A:  "m4a",
		FileFormatFLAC: "flac",
		FileFormatOPUS: "opus",
	}

	fileFormatByExtNoDot = map[string]FileFormat{
		"mp4":  FileFormatMP4,
		"webm": FileFormatWebM,
		"mp3":  FileFormatMP3,
		"m4a":  FileFormatM4A,
		"flac": FileFormatFLAC,
		"opus": FileFormatOPUS,
	}

	videoFormatByFileFormat = map[FileFormat]VideoFormat{
		FileFormatMP4:  VideoFormatMP4,
		FileFormatWebM: VideoFormatWebM,
	}

	audioFormatByFileFormat = map[FileFormat]AudioFormat{
		FileFormatMP3:  AudioFormatMP3,
		FileFormatM4A:  AudioFormatM4A,
		FileFormatFLAC: AudioFormatFLAC,
		FileFormatOPUS: AudioFormatOPUS,
	}
)

// String returns the value as a string.
func (v FileFormat) String() string {
	return string(v)
}

// ExtNoDot returns the file extension for the FileFormat, without a dot.
func (v FileFormat) ExtNoDot() string {
	return extNoDotByFileFormat[v]
}

// Ext returns the file extension for the FileFormat, prefixed with a dot.
func (v FileFormat) Ext() string {
	extNoDot := extNoDotByFileFormat[v]
	if extNoDot == "" {
		return ""
	}

	return "." + extNoDot
}

// Ptr returns the pointer.
func (v FileFormat) Ptr() *FileFormat {
	return &v
}

// Exists returns true if the FileFormat is valid.
func (v FileFormat) Exists() bool {
	_, exists := fileFormats[v]
	return exists
}

// IsVideo returns true if a video format.
func (v FileFormat) IsVideo() bool {
	_, exitsts := videoFormatByFileFormat[v]
	return exitsts
}

// IsAudio returns true if an audio format.
func (v FileFormat) IsAudio() bool {
	_, exitsts := audioFormatByFileFormat[v]
	return exitsts
}

// VideoFormat returns the corresponding VideoFormat for the FileFormat.
func (v FileFormat) VideoFormat() VideoFormat {
	videoFormat, exitsts := videoFormatByFileFormat[v]
	if !exitsts {
		return VideoFormatNone
	}
	return videoFormat
}

// AudioFormat returns the corresponding AudioFormat for the FileFormat.
func (v FileFormat) AudioFormat() AudioFormat {
	audioFormat, exitsts := audioFormatByFileFormat[v]
	if !exitsts {
		return AudioFormatNone
	}
	return audioFormat
}

// FormatType
func (v FileFormat) FormatType() FormatType {
	if v == FileFormatNone {
		return FormatTypeNone
	}

	if v.IsAudio() {
		return FormatTypeAudioOnly
	}

	return FormatTypeVideoAudio
}

// ParseFileFormat converting string to FileFormat
func ParseFileFormat(format string) (FileFormat, error) {
	if format == "" {
		return FileFormatNone, errors.New("empty value for FileFormat")
	}

	fileFormat := FileFormat(strings.ToLower(format))

	if !fileFormat.Exists() {
		return FileFormatNone, errors.New("invalid value for FileFormat")
	}

	return fileFormat, nil
}

// MustParseFileFormat converts a string to FileFormat, returning FileFormatNone if invalid.
func MustParseFileFormat(format string) FileFormat {
	fileFormat, err := ParseFileFormat(format)
	if err != nil {
		return FileFormatNone
	}

	return fileFormat
}

// FileExtToFileFormat converts a file extension to a FileFormat, returning FileFormatNone if not found.
func FileExtToFileFormat(ext string) FileFormat {
	extNoDot := strings.TrimPrefix(ext, ".")
	if extNoDot == "" {
		return FileFormatNone
	}

	fileFormat, exists := fileFormatByExtNoDot[strings.ToLower(extNoDot)]
	if !exists {
		return FileFormatNone
	}

	return fileFormat
}

// ValidateFileFormat checks if the field value is a valid FileFormat enum.
func ValidateFileFormat(fl validator.FieldLevel) bool {
	return FileFormat(fl.Field().String()).Exists()
}
