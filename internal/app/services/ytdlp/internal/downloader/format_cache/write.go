package formatcache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neosy/elengrab/pkg/nfile"
)

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
