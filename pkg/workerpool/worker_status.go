package workerpool

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Worker status
type WorkerStatus string

const (
	WorkerStatusNone    WorkerStatus = "none"
	WorkerStatusIdle    WorkerStatus = "idle"
	WorkerStatusWorking WorkerStatus = "working"
	WorkerStatusStopped WorkerStatus = "stopped"
)

var (
	// workerStatusMap implementation of a set for WorkerStatus
	workerStatusMap = map[WorkerStatus]struct{}{
		WorkerStatusNone:    {},
		WorkerStatusIdle:    {},
		WorkerStatusWorking: {},
	}
)

// String returns the value as a string.
func (v WorkerStatus) String() string {
	return string(v)
}

// Exists returns true if the WorkerStatus is valid.
func (v WorkerStatus) Exists() bool {
	_, exists := workerStatusMap[v]
	return exists
}

// ParseWorkerStatus converting string to status
func ParseWorkerStatus(s string) (WorkerStatus, error) {
	status := WorkerStatus(strings.ToUpper(s))

	if _, exists := workerStatusMap[status]; !exists {
		return "", errors.New("invalid value for WorkerStatus")
	}

	return status, nil
}

// StringToWorkerStatus converting string to status
func StringToWorkerStatus(s string) (WorkerStatus, error) {
	return ParseWorkerStatus(s)
}

// ValidateWorkerStatus checks if the field value is a valid WorkerStatus enum.
func ValidateWorkerStatus(fl validator.FieldLevel) bool {
	return WorkerStatus(fl.Field().String()).Exists()
}
