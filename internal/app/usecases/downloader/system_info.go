package downloader

import (
	"sync"

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

func (uc *Downloader) SystemInfo() dto.SystemInfoResponse {
	return uc.systemInfoStore.read()
}

func (uc *Downloader) UpdateSystemInfo() {
	stats, err := uc.downloadsStorage.Stats()
	if err != nil {
		uc.logger.Warn("Failed to get storage stats", "error", err)
	}

	systemInfoOld := uc.systemInfoStore.data

	uc.systemInfoStore.mu.Lock()
	uc.systemInfoStore.data.DiskUsed = stats.Used
	uc.systemInfoStore.data.DiskFree = stats.Free
	uc.systemInfoStore.mu.Unlock()

	if systemInfoOld.DiskFree != uc.systemInfoStore.data.DiskFree ||
		systemInfoOld.DiskUsed != uc.systemInfoStore.data.DiskUsed {
		uc.broadcastSystemInfoUpdate()
	}
}
