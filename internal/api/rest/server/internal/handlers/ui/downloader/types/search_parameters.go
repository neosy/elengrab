package types

import (
	"fmt"
	"strings"

	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

type SearchParameter struct {
	Key   qkeys.QueryKey
	Value string
}

type SearchParameters struct {
	items      []SearchParameter
	itemsByKey map[qkeys.QueryKey]SearchParameter
}

func NewSearchParameter(key qkeys.QueryKey, value string) SearchParameter {
	return SearchParameter{
		Key:   key,
		Value: value,
	}
}

func NewSearchParameters(parameters ...SearchParameter) *SearchParameters {
	sp := &SearchParameters{
		itemsByKey: make(map[qkeys.QueryKey]SearchParameter),
	}

	for _, p := range parameters {
		sp.Append(p)
	}

	return sp
}

func (p SearchParameter) QueryString() string {
	return fmt.Sprintf("%s=%s", p.Key.String(), p.Value)
}

func (p SearchParameter) QueryShortString() string {
	return fmt.Sprintf("%s=%s", p.Key.Short(), p.Value)
}

func (sp *SearchParameters) Append(parameter SearchParameter) *SearchParameters {
	sp.items = append(sp.items, parameter)
	sp.itemsByKey[parameter.Key] = parameter
	return sp
}

func (sp *SearchParameters) Add(key qkeys.QueryKey, value string) {
	sp.Append(NewSearchParameter(key, value))
}

func (sp *SearchParameters) QueryString() string {
	parameters := []string{}

	for _, p := range sp.items {
		parameters = append(parameters, p.QueryString())
	}

	return strings.Join(parameters, "&")
}

func (sp *SearchParameters) QueryShortString() string {
	parameters := []string{}

	for _, p := range sp.items {
		parameters = append(parameters, p.QueryShortString())
	}

	return strings.Join(parameters, "&")
}

func (sp *SearchParameters) QueryParameter() string {
	return fmt.Sprintf("%s=%s", qkeys.SearchParametersKey.String(), idcodec.EncodeStringBase64URL(sp.QueryShortString()))
}

func ParseEncodedSearchParameters(value string) (*SearchParameters, error) {
	queryString, err := idcodec.DecodeStringBase64URL(value)
	if err != nil {
		return nil, err
	}

	parameters := strings.Split(queryString, "&")
	if len(parameters) == 0 {
		return nil, nil
	}

	searchParameters := NewSearchParameters()

	for _, p := range parameters {
		parts := strings.SplitN(p, "=", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid search parameter: %s", p)
		}

		shortKey, value := parts[0], parts[1]

		key := qkeys.Keys.KeyByShortKey(shortKey)
		if key == "" {
			return nil, fmt.Errorf("unknown search parameter key: %s", shortKey)
		}

		searchParameters.Add(key, value)
	}

	return searchParameters, nil
}

func (sp *SearchParameters) ParseLastCursor() (*MediaQueryCursor, error) {
	lastCursorKey, exists := sp.itemsByKey[qkeys.LastCursorKey]
	if !exists {
		return &MediaQueryCursor{
			ViewMode: dtypes.QueryMediaViewModeDefault,
		}, nil
	}

	lastCursor, err := DecodeMediaQueryCursor(lastCursorKey.Value)
	if err != nil {
		return nil, err
	}

	return lastCursor, nil
}
