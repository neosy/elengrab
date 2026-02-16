package formatcache

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	formatCacheTTL = 2 * time.Hour
)

type FormatCache struct {
	mu  sync.RWMutex
	dir string
}

func NewFormatCache(dir string) *FormatCache {
	return &FormatCache{
		dir: dir,
	}
}

func (c *FormatCache) CacheDir() string {
	return c.dir
}

func (c *FormatCache) CacheFileName(url string) string {
	return FormatCacheFileName(url)
}

func (c *FormatCache) CacheFilePath(url string) string {
	filePath := filepath.Join(c.dir, c.CacheFileName(url))
	return filePath
}

func (c *FormatCache) IsTTLValidByURL(url string) (bool, error) {
	return c.IsTTLValid(c.CacheFilePath(url))
}

func (c *FormatCache) IsTTLValid(filePath string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	expiredTime := fileInfo.ModTime().Add(formatCacheTTL)
	return time.Now().UTC().Before(expiredTime), nil
}
