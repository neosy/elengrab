package executor

import (
	"context"
	"fmt"
	"os"
)

func (e *Executor) EnsureFormatCache(
	ctx context.Context,
	url string,
	useCookies bool,
) error {
	valid, err := e.formatCache.IsTTLValidByURL(url)
	if err != nil {
		if os.IsNotExist(err) {
			e.logger.Debug("File format cache miss", "url", url)
		} else {
			return fmt.Errorf("failed to check TTL validity for URL %q: %w", url, err)
		}
	}

	if valid {
		e.logger.Debug(
			"Format cache TTL is still valid, no update needed",
			"url", url,
		)
		return nil
	}

	if err == nil {
		e.logger.Debug("Format cache TTL expired", "url", url)
	}

	_, err = e.fetchAndCacheInfoJSON(ctx, url, useCookies)
	if err != nil {
		return fmt.Errorf("failed to fetch and cache formats JSON for URL %q: %w", url, err)
	}

	return nil
}
