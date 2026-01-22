package ddownload

// Progress of the download
type DownloadProgress struct {
	DownloadedBytes  int64
	TotalBytes       int64
	ETASeconds       int
	SpeedBytesPerSec int64
}

// Percent returns the download progress percentage.
func (p DownloadProgress) Percent() float64 {
	if p.TotalBytes == 0 {
		return 0
	}
	return float64(p.DownloadedBytes) / float64(p.TotalBytes) * 100
}
