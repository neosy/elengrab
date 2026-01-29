package idto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/neosy/elengrab/pkg/nfile"
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

func (c *FormatCache) WriteByURL(url string, data []byte) error {
	return c.write(c.CacheFilePath(url), data)
}

func (c *FormatCache) write(filePath string, data []byte) error {
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

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	_ = os.Remove(filePath)
	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func (c *FormatCache) delete(filePath string) error {
	return os.Remove(filePath)
}

func (c *FormatCache) DeleteByURL(url string) error {
	return c.delete(c.CacheFilePath(url))
}

func (c *FormatCache) LoadByURL(url string) ([]byte, error) {
	return c.load(c.CacheFilePath(url))
}

func (c *FormatCache) load(filePath string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	isTTLValid, err := c.IsTTLValid(filePath)
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
	return time.Now().Before(expiredTime), nil
}

func FormatCacheFileName(url string) string {
	hash := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%s.json", hex.EncodeToString(hash[:16]))
}
