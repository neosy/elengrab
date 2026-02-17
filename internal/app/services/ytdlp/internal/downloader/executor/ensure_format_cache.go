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
		return nil
	}

	_, err = e.fetchAndCacheFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return fmt.Errorf("failed to fetch and cache formats JSON for URL %q: %w", url, err)
	}

	return nil
}
