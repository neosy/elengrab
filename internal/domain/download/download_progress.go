package ddownload

import uptr "github.com/neosy/elengrab/pkg/utils/pointer"

// Progress of the download
type DownloadProgress struct {
	// Total bytes downloaded
	DownloadedBytes int64
	// Total bytes to download
	TotalBytes int64
	// ETA Seconds remaining
	ETASeconds int
	// Speed bytes per second
	SpeedBytesPerSec int64
}

// Percent returns the download progress percentage.
func (p DownloadProgress) Percent() float64 {
	if p.TotalBytes == 0 {
		return 0
	}
	return float64(p.DownloadedBytes) / float64(p.TotalBytes) * 100
}

func (src *DownloadProgress) Copy() *DownloadProgress {
	return uptr.Copy(src)
}
