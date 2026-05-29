package dltypes

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type WorkingStatus string

const (
	WorkingStatusNone           WorkingStatus = "none"
	WorkingStatusStartDownload  WorkingStatus = "start"
	WorkingStatusDownloading    WorkingStatus = "downloading"
	WorkingStatusFinishDownload WorkingStatus = "finish"
)

var (
	// workingStatusMap implementation of a set for WorkingStatus
	workingStatusMap = map[WorkingStatus]struct{}{
		WorkingStatusNone:           {},
		WorkingStatusStartDownload:  {},
		WorkingStatusDownloading:    {},
		WorkingStatusFinishDownload: {},
	}

	usecaseWorkingStatusToUI = map[dto.WorkingStatus]WorkingStatus{
		dto.WorkingStatusNone:           WorkingStatusNone,
		dto.WorkingStatusStartDownload:  WorkingStatusStartDownload,
		dto.WorkingStatusDownloading:    WorkingStatusDownloading,
		dto.WorkingStatusFinishDownload: WorkingStatusFinishDownload,
	}
)

// String returns the value as a string.
func (v WorkingStatus) String() string {
	return string(v)
}

// Exists returns true if the WorkingStatus is valid.
func (v WorkingStatus) Exists() bool {
	_, exists := workingStatusMap[v]
	return exists
}

func MapUsecaseWorkingStatusToUI(workingStatus dto.WorkingStatus) WorkingStatus {
	status, ok := usecaseWorkingStatusToUI[workingStatus]
	if !ok {
		return WorkingStatusNone
	}
	return status
}
