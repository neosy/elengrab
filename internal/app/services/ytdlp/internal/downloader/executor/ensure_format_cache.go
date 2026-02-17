package executor

import (
	"context"
	"fmt"
)

func (e *Executor) EnsureFormatCache(
	ctx context.Context,
	url string,
	useCookies bool,
) error {
	valid, err := e.formatCache.IsTTLValidByURL(url)
	if err != nil {
		return fmt.Errorf("failed to check TTL validity for URL %q: %w", url, err)
	}

	if valid {
		e.logger.Debug(
			"Format cache TTL is still valid, no update needed",
			"url", url,
		)
		return nil
	}

	e.logger.Debug("Format cache TTL expired", "url", url)

	_, err = e.fetchAndCacheFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return fmt.Errorf("failed to fetch and cache formats JSON for URL %q: %w", url, err)
	}

	return nil
}
