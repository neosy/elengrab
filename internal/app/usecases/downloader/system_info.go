package downloader

import (
	"sync"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/helper"
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
	used := helper.FolderSize(uc.downloadsDir)
	free, _, _, _ := helper.DiskUsage(uc.downloadsDir)

	systemInfoOld := uc.systemInfoStore.data

	uc.systemInfoStore.mu.Lock()
	uc.systemInfoStore.data.DiskUsed = used
	uc.systemInfoStore.data.DiskFree = int64(free)
	uc.systemInfoStore.mu.Unlock()

	if systemInfoOld.DiskFree != uc.systemInfoStore.data.DiskFree ||
		systemInfoOld.DiskUsed != uc.systemInfoStore.data.DiskUsed {
		uc.broadcastSystemInfoUpdate()
	}
}
