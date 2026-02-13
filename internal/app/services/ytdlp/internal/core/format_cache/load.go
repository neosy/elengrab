package formatcache

import "os"

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
