package dto

// Working status
type WorkingStatus string

const (
	WorkingStatusNone           WorkingStatus = "none"
	WorkingStatusStartDownload  WorkingStatus = "start_download"
	WorkingStatusDownloading    WorkingStatus = "downloading"
	WorkingStatusFinishDownload WorkingStatus = "finish_download"
)

var (
	// workingStatusMap implementation of a set for WorkingStatus
	workingStatusMap = map[WorkingStatus]struct{}{
		WorkingStatusNone:           {},
		WorkingStatusStartDownload:  {},
		WorkingStatusDownloading:    {},
		WorkingStatusFinishDownload: {},
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
