package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// MediaDownload status
type MediaDownloadStatus uint8

const (
	MediaDownloadStatusNone MediaDownloadStatus = iota
	MediaDownloadStatusNew
	MediaDownloadStatusPending
	MediaDownloadStatusWorking
	MediaDownloadStatusDone
	MediaDownloadStatusFailed
	MediaDownloadStatusRefreshing

	MediaDownloadStatusDefault = MediaDownloadStatusNone
)

var (
	// mediaDownloadStatusNameByStatus implementation of a set for MediaDownloadStatus
	mediaDownloadStatusNameByStatus = map[MediaDownloadStatus]string{
		MediaDownloadStatusNone:       "none",
		MediaDownloadStatusNew:        "new",
		MediaDownloadStatusPending:    "pending",
		MediaDownloadStatusWorking:    "working",
		MediaDownloadStatusDone:       "done",
		MediaDownloadStatusFailed:     "failed",
		MediaDownloadStatusRefreshing: "refreshing",
	}

	parseMediaDownloadStatusMap = map[string]MediaDownloadStatus{
		"none":       MediaDownloadStatusNone,
		"new":        MediaDownloadStatusNew,
		"pending":    MediaDownloadStatusPending,
		"working":    MediaDownloadStatusWorking,
		"done":       MediaDownloadStatusDone,
		"failed":     MediaDownloadStatusFailed,
		"refreshing": MediaDownloadStatusRefreshing,
	}
)

// String returns the value as a string.
func (v MediaDownloadStatus) String() string {
	return mediaDownloadStatusNameByStatus[v]
}

// IsReady returns true when the download status indicates that the media is ready for use.
func (v MediaDownloadStatus) IsReady() bool {
	return v == MediaDownloadStatusDone ||
		v == MediaDownloadStatusRefreshing
}

// Exists returns true if the MediaDownloadStatus is valid.
func (v MediaDownloadStatus) Exists() bool {
	_, exists := mediaDownloadStatusNameByStatus[v]
	return exists
}

// ParseMediaDownloadStatus converting string to MediaDownloadStatus
func ParseMediaDownloadStatus(s string) (MediaDownloadStatus, error) {
	status, exists := parseMediaDownloadStatusMap[strings.ToLower(s)]
	if !exists {
		return MediaDownloadStatusDefault, errors.New("invalid value for MediaDownloadStatus")
	}
	return status, nil
}

// MustParseMediaDownloadStatus converting string to MediaDownloadStatus, ignoring any errors.
func MustParseMediaDownloadStatus(s string) MediaDownloadStatus {
	status, _ := ParseMediaDownloadStatus(s)
	return status
}

// ValidateMediaDownloadStatus checks if the field value is a valid MediaDownloadStatus enum.
func ValidateMediaDownloadStatus(fl validator.FieldLevel) bool {
	_, exist := parseMediaDownloadStatusMap[fl.Field().String()]
	return exist
}
