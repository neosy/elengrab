package executor

import "context"

func (e *Executor) UpdateFormatCache(
	ctx context.Context,
	url string,
	useCookies bool,
) error {
	valid, err := e.formatCache.IsTTLValidByURL(url)
	if err != nil {
		return err
	}

	if valid {
		return nil
	}

	return e.updateFormatCacheForce(ctx, url, useCookies)
}

func (e *Executor) updateFormatCacheForce(
	ctx context.Context,
	url string,
	useCookies bool,
) error {
	_, err := e.fetchAndCacheFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return err
	}

	return nil
}
