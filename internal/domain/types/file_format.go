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
	fileFormatMap = map[FileFormat]struct{}{
		FileFormatNone: {},
		FileFormatAuto: {},
		FileFormatMP4:  {},
		FileFormatWebM: {},
		FileFormatMP3:  {},
		FileFormatM4A:  {},
		FileFormatFLAC: {},
		FileFormatOPUS: {},
	}

	fileExtToFileFormatMap = map[string]FileFormat{
		"mp4":  FileFormatMP4,
		"webm": FileFormatWebM,
		"mp3":  FileFormatMP3,
		"m4a":  FileFormatM4A,
		"flac": FileFormatFLAC,
		"opus": FileFormatOPUS,
	}
)

// String returns the value as a string.
func (v FileFormat) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v FileFormat) Ptr() *FileFormat {
	return &v
}

// Exists returns true if the FileFormat is valid.
func (v FileFormat) Exists() bool {
	_, exists := fileFormatMap[v]
	return exists
}

// ParseFileFormat converting string to FileFormat
func ParseFileFormat(s string) (FileFormat, error) {
	fileFormat := FileFormat(strings.ToLower(s))

	if _, exists := fileFormatMap[fileFormat]; !exists {
		return "", errors.New("invalid value for FileFormat")
	}

	return fileFormat, nil
}

// MapFileExtToFileFormat
func MapFileExtToFileFormat(ext string) FileFormat {
	format, exist := fileExtToFileFormatMap[ext]
	if !exist {
		format = FileFormatNone
	}
	return format
}

// ValidateFileFormat checks if the field value is a valid FileFormat enum.
func ValidateFileFormat(fl validator.FieldLevel) bool {
	return FileFormat(fl.Field().String()).Exists()
}
