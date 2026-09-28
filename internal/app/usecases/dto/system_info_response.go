package dto

type SystemInfoResponse struct {
	AppVersion string `json:"appVersion"`

	GoVersion     string `json:"goVersion"`
	YtDlpVersion  string `json:"ytDlpVersion"`
	FFmpegVersion string `json:"ffmpegVersion"`
	DenoVersion   string `json:"denoVersion"`

	DiskFree uint64 `json:"diskFree"`
	DiskUsed uint64 `json:"diskUsed"`
}
