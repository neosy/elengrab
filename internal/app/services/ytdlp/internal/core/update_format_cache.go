package core

import "context"

func (c *Core) updateFormatCache(
	ctx context.Context,
	url string,
	useCookies bool,
) error {
	valid, err := c.formatCache.IsTTLValidByURL(url)
	if err != nil {
		return err
	}

	if valid {
		return nil
	}

	return c.updateFormatCacheForce(ctx, url, useCookies)
}

func (c *Core) updateFormatCacheForce(
	ctx context.Context,
	url string,
	useCookies bool,
) error {
	_, err := c.fetchAndCacheFormatsJSON(ctx, url, useCookies)
	if err != nil {
		return err
	}

	return nil
}
