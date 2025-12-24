package ytdlp

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/internal/yt_dlp/helper"
	"github.com/neosy/elengrab/pkg/nfile"
)

const (
	formatCacheTTL = 2 * time.Hour
)

type formatCache struct {
	mu sync.RWMutex

	dir string
}

func NewFormatCache(dir string) *formatCache {
	return &formatCache{
		dir: dir,
	}
}

func (c *formatCache) cacheDir() string {
	return c.dir
}

func (c *formatCache) cacheFileName(url string) string {
	return helper.FormatCacheFileName(url)
}

func (c *formatCache) cacheFilePath(url string) string {
	filePath := filepath.Join(c.dir, c.cacheFileName(url))
	return filePath
}

func (c *formatCache) writeByURL(url string, data []byte) error {
	return c.write(c.cacheFilePath(url), data)
}

func (c *formatCache) write(filePath string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure the cache directory exists
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Full path to the cache file
	var (
		tmpFilePattern = fmt.Sprintf(
			"%s-*%s",
			nfile.FileNameWithoutExt(filepath.Base(filePath)),
			filepath.Ext(filePath),
		)
	)

	tmpFile, err := os.CreateTemp(c.dir, tmpFilePattern)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer tmpFile.Close()

	_, err = tmpFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	_ = os.Remove(filePath)
	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func (c *formatCache) delete(filePath string) error {
	return os.Remove(filePath)
}

func (c *formatCache) deleteByURL(url string) error {
	return c.delete(c.cacheFilePath(url))
}

func (c *formatCache) loadByURL(url string) ([]byte, error) {
	return c.load(c.cacheFilePath(url))
}

func (c *formatCache) load(filePath string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	isTTLValid, err := c.isTTLValid(filePath)
	if err != nil {
		return nil, err
	}

	if !isTTLValid {
		return nil, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

func (c *formatCache) isTTLValidByURL(url string) (bool, error) {
	return c.isTTLValid(c.cacheFilePath(url))
}

func (c *formatCache) isTTLValid(filePath string) (bool, error) {
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
	return time.Now().Before(expiredTime), nil
}
