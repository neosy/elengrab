package formatcache

import "os"

func (c *FormatCache) delete(filePath string) error {
	return os.Remove(filePath)
}

func (c *FormatCache) DeleteByURL(url string) error {
	return c.delete(c.CacheFilePath(url))
}
