package httppaths

import (
	"net/url"
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func buildImageSufix(verHash string, sources []dtypes.ImageSource) string {
	var sourceStrings []string
	for _, src := range sources {
		if src.Exists() {
			sourceStrings = append(sourceStrings, src.String())
		}
	}

	var urlValues url.Values
	if len(sourceStrings) > 0 {
		urlValues = url.Values{}
		urlValues.Set("source", strings.Join(sourceStrings, ","))
	}

	if verHash != "" {
		urlValues.Set("v", verHash)
	}

	var urlSufix string
	if len(urlValues) > 0 {
		urlSufix = "?" + urlValues.Encode()
	}

	return urlSufix
}
