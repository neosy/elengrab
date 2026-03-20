package dto

type SystemInfoResponse struct {
	AppVersion string `json:"appVersion"`
	DiskFree   int64  `json:"diskFree"`
	DiskUsed   int64  `json:"diskUsed"`
}
