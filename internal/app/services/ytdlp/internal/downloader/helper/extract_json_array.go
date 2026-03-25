package helper

import (
	"fmt"
	"strings"

	"github.com/neosy/elengrab/internal/pkg/errorx"
)

// ExtractJSONArray searches for a JSON array in the HTML starting from the given key.
// It balances the square brackets to ensure the extracted string is a valid JSON array.
//
//	exemple key: `"decoratedAvatarViewModel":{"avatar":{"avatarViewModel":{"image":{"sources":[`
//	include the opening bracket '['
func ExtractJSONArray(html, key string) (string, error) {
	start := strings.Index(html, key)
	if start == -1 {
		return "", errorx.Errorf("json: key %q not found", key)
	}
	start += len(key) - 1 // include the opening bracket '['

	brackets := 1
	end := start + 1
	for end < len(html) && brackets > 0 {
		switch html[end] {
		case '[':
			brackets++
		case ']':
			brackets--
		}
		end++
	}

	if brackets != 0 {
		return "", fmt.Errorf("unbalanced brackets in JSON array")
	}

	return html[start:end], nil
}
