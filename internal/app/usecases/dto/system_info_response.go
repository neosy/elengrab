package dto

type SystemInfoResponse struct {
	AppVersion string `json:"appVersion"`
	DiskFree   uint64 `json:"diskFree"`
	DiskUsed   uint64 `json:"diskUsed"`
}
