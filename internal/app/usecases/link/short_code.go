package link

import (
	"github.com/google/uuid"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (u *Link) GenerateShortCode(url string, length uint8, deterministic bool) string {
	return dlink.GenerateShortCode(uuid.Nil, url, length, true)
}
