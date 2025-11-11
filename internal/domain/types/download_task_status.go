package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Download task status
type DownloadTaskStatus string

const (
	DownloadTaskStatusNone    DownloadTaskStatus = "none"
	DownloadTaskStatusPending DownloadTaskStatus = "pending"
	DownloadTaskStatusWorking DownloadTaskStatus = "working"
)

var (
	// downloadTaskStatusMap implementation of a set for DownloadTaskStatus
	downloadTaskStatusMap = map[DownloadTaskStatus]struct{}{
		DownloadTaskStatusNone:    {},
		DownloadTaskStatusPending: {},
		DownloadTaskStatusWorking: {},
	}
)

// String returns the value as a string.
func (v DownloadTaskStatus) String() string {
	return string(v)
}

// Exists returns true if the DownloadTaskStatus is valid.
func (v DownloadTaskStatus) Exists() bool {
	_, exists := downloadTaskStatusMap[v]
	return exists
}

// ParseDownloadTaskStatus converting string to order status
func ParseDownloadTaskStatus(s string) (DownloadTaskStatus, error) {
	status := DownloadTaskStatus(strings.ToUpper(s))

	if _, exists := downloadTaskStatusMap[status]; !exists {
		return "", errors.New("invalid value for DownloadTaskStatus")
	}

	return status, nil
}

// StringToDownloadTaskStatus converting string to order status
func StringToDownloadTaskStatus(s string) (DownloadTaskStatus, error) {
	return ParseDownloadTaskStatus(s)
}

// ValidateDownloadTaskStatus checks if the field value is a valid DownloadTaskStatus enum.
func ValidateDownloadTaskStatus(fl validator.FieldLevel) bool {
	return DownloadTaskStatus(fl.Field().String()).Exists()
}
