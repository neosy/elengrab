package downloader

import (
	"runtime"
	"strings"
	"sync"

	denosrv "github.com/neosy/elengrab/internal/app/services/deno"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

type systemInfoStore struct {
	mu   sync.RWMutex
	data dto.SystemInfoResponse
}

func (s *systemInfoStore) read() dto.SystemInfoResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}

func (uc *downloader) SystemInfo() dto.SystemInfoResponse {
	return uc.systemInfoStore.read()
}

func (uc *downloader) UpdateSystemDiskInfo() {
	stats, err := uc.downloadsStorage.Stats()
	if err != nil {
		uc.logger.Warn("Failed to get storage stats", "error", err)
	}

	used, err := uc.downloadsStorage.Used()
	if err != nil {
		uc.logger.Warn("Failed to get used storage", "error", err)
	}

	systemInfoOld := uc.systemInfoStore.data

	uc.systemInfoStore.mu.Lock()
	uc.systemInfoStore.data.DiskUsed = used
	uc.systemInfoStore.data.DiskFree = stats.Free
	uc.systemInfoStore.mu.Unlock()

	if systemInfoOld.DiskFree != uc.systemInfoStore.data.DiskFree ||
		systemInfoOld.DiskUsed != uc.systemInfoStore.data.DiskUsed {
		uc.broadcastSystemInfoUpdate()
	}
}
func (uc *downloader) UpdateSystemVersionInfo() {
	ytDlpVersion, _ := uc.downloaderSrv.GetVersion(uc.appCtx)
	ffmpegVersion, _ := uc.ffmpegSrv.GetFFmpegVersion(uc.appCtx)

	denoVersion, _ := denosrv.GetVersion(uc.appCtx)

	uc.systemInfoStore.mu.Lock()
	uc.systemInfoStore.data.GoVersion = strings.TrimPrefix(runtime.Version(), "go")
	uc.systemInfoStore.data.YtDlpVersion = ytDlpVersion
	uc.systemInfoStore.data.FFmpegVersion = ffmpegVersion
	uc.systemInfoStore.data.DenoVersion = denoVersion
	uc.systemInfoStore.mu.Unlock()
}

func (uc *downloader) UpdateSystemInfo() {
	uc.UpdateSystemDiskInfo()
	uc.UpdateSystemVersionInfo()
}
