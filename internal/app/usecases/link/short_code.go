package link

import (
	"github.com/google/uuid"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (u *Link) GenerateShortCodeByURL(url string) string {
	return dlink.GenerateShortCode(uuid.Nil, url, u.options.ShortCodeLength, true)
}
