package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// File status
type FileStatus string

const (
	FileStatusNone    FileStatus = "none"
	FileStatusNew     FileStatus = "new"
	FileStatusPending FileStatus = "pending"
	FileStatusWorking FileStatus = "working"
	FileStatusDone    FileStatus = "done"
	FileStatusFailed  FileStatus = "failed"
)

var (
	// fileStatusMap implementation of a set for FileStatus
	fileStatusMap = map[FileStatus]struct{}{
		FileStatusNone:    {},
		FileStatusNew:     {},
		FileStatusPending: {},
		FileStatusWorking: {},
		FileStatusDone:    {},
		FileStatusFailed:  {},
	}
)

// String returns the value as a string.
func (v FileStatus) String() string {
	return string(v)
}

// Exists returns true if the FileStatus is valid.
func (v FileStatus) Exists() bool {
	_, exists := fileStatusMap[v]
	return exists
}

// ParseFileStatus converting string to status
func ParseFileStatus(s string) (FileStatus, error) {
	status := FileStatus(strings.ToUpper(s))

	if _, exists := fileStatusMap[status]; !exists {
		return "", errors.New("invalid value for FileStatus")
	}

	return status, nil
}

// StringToFileStatus converting string to status
func StringToFileStatus(s string) (FileStatus, error) {
	return ParseFileStatus(s)
}

// ValidateFileStatus checks if the field value is a valid FileStatus enum.
func ValidateFileStatus(fl validator.FieldLevel) bool {
	return FileStatus(fl.Field().String()).Exists()
}
