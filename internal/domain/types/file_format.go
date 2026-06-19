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

	fileFormatToExtMap = map[FileFormat]string{
		FileFormatMP4:  "mp4",
		FileFormatWebM: "webm",
		FileFormatMP3:  "mp3",
		FileFormatM4A:  "m4a",
		FileFormatFLAC: "flac",
		FileFormatOPUS: "opus",
	}

	fileExtToFileFormatMap = map[string]FileFormat{
		"mp4":  FileFormatMP4,
		"webm": FileFormatWebM,
		"mp3":  FileFormatMP3,
		"m4a":  FileFormatM4A,
		"flac": FileFormatFLAC,
		"opus": FileFormatOPUS,
	}

	fileFormatToVideoFormatMap = map[FileFormat]VideoFormat{
		FileFormatMP4:  VideoFormatMP4,
		FileFormatWebM: VideoFormatWebM,
	}

	fileFormatToAudioFormatMap = map[FileFormat]AudioFormat{
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

// Ext returns the file extension for the FileFormat.
func (v FileFormat) Ext() string {
	return fileFormatToExtMap[v]
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

// IsVideo returns true if a video format.
func (v FileFormat) IsVideo() bool {
	_, exitsts := fileFormatToVideoFormatMap[v]
	return exitsts
}

// IsAudio returns true if an audio format.
func (v FileFormat) IsAudio() bool {
	_, exitsts := fileFormatToAudioFormatMap[v]
	return exitsts
}

// VideoFormat returns the corresponding VideoFormat for the FileFormat.
func (v FileFormat) VideoFormat() VideoFormat {
	videoFormat, exitsts := fileFormatToVideoFormatMap[v]
	if !exitsts {
		return VideoFormatNone
	}
	return videoFormat
}

// AudioFormat returns the corresponding AudioFormat for the FileFormat.
func (v FileFormat) AudioFormat() AudioFormat {
	audioFormat, exitsts := fileFormatToAudioFormatMap[v]
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
	fileFormat := FileFormat(strings.ToLower(format))

	if _, exists := fileFormatMap[fileFormat]; !exists {
		return "", errors.New("invalid value for FileFormat")
	}

	return fileFormat, nil
}

// MapFileExtToFileFormat
func MapFileExtToFileFormat(ext string) FileFormat {
	ext, _ = strings.CutPrefix(ext, ".")
	ext = strings.ToLower(ext)

	format, exist := fileExtToFileFormatMap[ext]
	if !exist {
		return FileFormatNone
	}

	return format
}

// ValidateFileFormat checks if the field value is a valid FileFormat enum.
func ValidateFileFormat(fl validator.FieldLevel) bool {
	return FileFormat(fl.Field().String()).Exists()
}
