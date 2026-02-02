package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Download task status
type DownloadTaskStatus uint8

const (
	DownloadTaskStatusNone DownloadTaskStatus = iota
	DownloadTaskStatusNew
	DownloadTaskStatusPending
	DownloadTaskStatusWorking
	DownloadTaskStatusFailed

	DownloadTaskStatusDefault = DownloadTaskStatusNone
)

var (
	// downloadTaskStatusMap implementation of a set for DownloadTaskStatus
	downloadTaskStatusMap = map[DownloadTaskStatus]string{
		DownloadTaskStatusNone:    "none",
		DownloadTaskStatusNew:     "new",
		DownloadTaskStatusPending: "pending",
		DownloadTaskStatusWorking: "working",
		DownloadTaskStatusFailed:  "failed",
	}

	parserDownloadTaskStatusMap = map[string]DownloadTaskStatus{
		"none":    DownloadTaskStatusNone,
		"new":     DownloadTaskStatusNew,
		"pending": DownloadTaskStatusPending,
		"working": DownloadTaskStatusWorking,
		"failed":  DownloadTaskStatusFailed,
	}
)

// String returns the value as a string.
func (v DownloadTaskStatus) String() string {
	return downloadTaskStatusMap[v]
}

// Exists returns true if the DownloadTaskStatus is valid.
func (v DownloadTaskStatus) Exists() bool {
	_, exists := downloadTaskStatusMap[v]
	return exists
}

// ParseDownloadTaskStatus converting string to status
func ParseDownloadTaskStatus(s string) (DownloadTaskStatus, error) {
	status, exists := parserDownloadTaskStatusMap[strings.ToLower(s)]
	if !exists {
		return DownloadTaskStatusDefault, errors.New("invalid value for DownloadTaskStatus")
	}
	return status, nil
}

// MustParseDownloadTaskStatus converting string to DownloadTaskStatus, ignoring any errors.
func MustParseDownloadTaskStatus(s string) DownloadTaskStatus {
	status, _ := ParseDownloadTaskStatus(s)
	return status
}

// ValidateDownloadTaskStatus checks if the field value is a valid DownloadTaskStatus enum.
func ValidateDownloadTaskStatus(fl validator.FieldLevel) bool {
	_, exist := parserDownloadTaskStatusMap[fl.Field().String()]
	return exist
}
