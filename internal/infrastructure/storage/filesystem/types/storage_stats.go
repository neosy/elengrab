package storagetypes

// StorageStats represents disk usage information.
type StorageStats struct {
	Total uint64
	Used  uint64
	Free  uint64
}
