package linkweb

import (
	"context"
)

func (u *LinkWeb) DeleteShortLink(ctx context.Context, url string) error {
	shortCode := u.link.GenerateShortCodeByURL(url, u.options.ShortCodeLength, true)

	link, err := u.link.GetLastByShortCode(ctx, shortCode)
	if err != nil {
		return err
	}

	return u.link.SoftDelete(ctx, link.LinkID)
}
