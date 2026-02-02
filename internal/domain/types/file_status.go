package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// File status
type FileStatus uint8

const (
	FileStatusNone FileStatus = iota
	FileStatusNew
	FileStatusPending
	FileStatusWorking
	FileStatusDone
	FileStatusFailed

	FileStatusDefault = FileStatusNone
)

var (
	// fileStatusMap implementation of a set for FileStatus
	fileStatusMap = map[FileStatus]string{
		FileStatusNone:    "none",
		FileStatusNew:     "new",
		FileStatusPending: "pending",
		FileStatusWorking: "working",
		FileStatusDone:    "done",
		FileStatusFailed:  "failed",
	}

	parseFileStatusMap = map[string]FileStatus{
		"none":    FileStatusNone,
		"new":     FileStatusNew,
		"pending": FileStatusPending,
		"working": FileStatusWorking,
		"done":    FileStatusDone,
		"failed":  FileStatusFailed,
	}
)

// String returns the value as a string.
func (v FileStatus) String() string {
	return fileStatusMap[v]
}

// Exists returns true if the FileStatus is valid.
func (v FileStatus) Exists() bool {
	_, exists := fileStatusMap[v]
	return exists
}

// ParseFileStatus converting string to FileStatus
func ParseFileStatus(s string) (FileStatus, error) {
	status, exists := parseFileStatusMap[strings.ToLower(s)]
	if !exists {
		return FileStatusDefault, errors.New("invalid value for FileStatus")
	}
	return status, nil
}

// MustParseFileStatus converting string to FileStatus, ignoring any errors.
func MustParseFileStatus(s string) FileStatus {
	status, _ := ParseFileStatus(s)
	return status
}

// ValidateFileStatus checks if the field value is a valid FileStatus enum.
func ValidateFileStatus(fl validator.FieldLevel) bool {
	_, exist := parseFileStatusMap[fl.Field().String()]
	return exist
}
