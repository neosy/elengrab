package link

import (
	"uuid"

	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (u *link) GenerateShortCode(url string, length uint8, deterministic bool) string {
	return dlink.GenerateShortCode(uuid.Nil(), url, length, true)
}
