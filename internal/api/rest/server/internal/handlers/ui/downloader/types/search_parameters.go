package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
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

func (p SearchParameter) ShortQueryString() string {
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

func (sp *SearchParameters) AddFromParmValues(values SearchParameterValues) {
	viewMode := dtypes.QueryMediaViewModeDefault.String()

	if values.ViewMode != "" {
		viewMode = values.ViewMode
	}

	lastCursor := MediaQueryCursor(values.LastRecord)

	sp.Add(qkeys.ViewModeKey, viewMode)
	if !lastCursor.IsZero() {
		sp.Add(qkeys.LastCursorKey, lastCursor.Encode())
	}

	if values.Filters != nil {
		for _, filter := range values.Filters.byKey {
			sp.Add(filter.Key, filter.Value)
		}
	}
}

func (sp *SearchParameters) AddValues(
	viewMode string,
	filters *QueryFilters,
	lastCursor dtypes.QueryMediaDownloadCursor,
) {
	sp.AddFromParmValues(SearchParameterValues{
		ViewMode:   viewMode,
		Filters:    filters,
		LastRecord: lastCursor,
	})
}

func (sp *SearchParameters) Find(key qkeys.QueryKey) (string, bool) {
	item, exists := sp.itemsByKey[key]
	if !exists {
		return "", false
	}

	return item.Value, true
}

func (sp *SearchParameters) FindViewMode() (dtypes.QueryMediaViewMode, error) {
	viewModeStr, exists := sp.Find(qkeys.ViewModeKey)
	if !exists {
		return dtypes.QueryMediaViewModeNone, nil
	}

	if viewModeStr == "" {
		return dtypes.QueryMediaViewModeNone, nil
	}

	mode, err := dtypes.ParseQueryMediaViewMode(viewModeStr)
	if err != nil {
		return dtypes.QueryMediaViewModeNone, err
	}

	return mode, nil
}

func (sp *SearchParameters) FindChannelID() uuid.UUID {
	channelID, exists := sp.Find(qkeys.ChannelIDKey)
	if !exists || channelID == "" {
		return uuid.Nil
	}

	id, err := idcodec.DecodeUUIDBase64URL(channelID)
	if err != nil {
		return uuid.Nil
	}

	return id
}

func (sp *SearchParameters) QueryString() string {
	parameters := []string{}

	for _, p := range sp.items {
		parameters = append(parameters, p.QueryString())
	}

	return strings.Join(parameters, "&")
}

func (sp *SearchParameters) ShortQueryString() string {
	parameters := []string{}

	for _, p := range sp.items {
		parameters = append(parameters, p.ShortQueryString())
	}

	return strings.Join(parameters, "&")
}

// EncodeShortQueryValue encodes the short query value using Base64URL.
func (sp *SearchParameters) EncodeShortQueryValue() string {
	return idcodec.EncodeStringBase64URL(sp.ShortQueryString())
}

// EncodeShortQueryString returns the short query as a key=value string.
func (sp *SearchParameters) EncodeShortQueryString() string {
	return fmt.Sprintf(
		"%s=%s",
		qkeys.SearchParametersKey.String(),
		sp.EncodeShortQueryValue(),
	)
}

func ParseEncodedSearchParamQueryString(queryString string) (*SearchParameters, error) {
	queryString, err := idcodec.DecodeStringBase64URL(queryString)
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

		key := qkeys.Keys.FindByShortKey(shortKey)
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
		return nil, nil
	}

	lastCursor, err := DecodeMediaQueryCursor(lastCursorKey.Value)
	if err != nil {
		return nil, err
	}

	return lastCursor, nil
}

func (sp *SearchParameters) ParseValues() (SearchParameterValues, error) {
	values := SearchParameterValues{
		ViewMode:   dtypes.QueryMediaViewModeDefault.String(),
		LastRecord: dtypes.QueryMediaDownloadCursor{},
	}

	if sp == nil {
		return values, nil
	}

	mode, err := sp.FindViewMode()
	if err != nil {
		return values, err
	}

	if mode != dtypes.QueryMediaViewModeNone {
		values.ViewMode = mode.String()
	}

	lastQueryCursor, err := sp.ParseLastCursor()
	if err != nil {
		return values, err
	}

	if lastQueryCursor != nil {
		values.LastRecord = dtypes.QueryMediaDownloadCursor(*lastQueryCursor)
	}

	values.Filters = sp.QueryFilters()

	return values, nil
}

func (sp *SearchParameters) QueryFilters() *QueryFilters {
	if len(sp.itemsByKey) == 0 {
		return nil
	}

	var filters *QueryFilters

	for key, param := range sp.itemsByKey {
		filterName := qkeys.Keys.FilterNameByKey(key)
		if filterName == dtypes.QueryFilterNameNone {
			continue
		}

		if filters == nil {
			filters = NewQueryFilters()
		}

		filters.Add(key, param.Value)
	}

	return filters
}

func (sp *SearchParameters) BuildJSON() []byte {
	if sp == nil {
		return nil
	}

	valuesByKey := make(map[string]string)

	for _, param := range sp.itemsByKey {
		valuesByKey[param.Key.String()] = param.Value
	}

	json, err := json.Marshal(valuesByKey)
	if err != nil {
		return nil
	}

	return json
}

func ParseSearchQueryString(queryString string) (*SearchParameters, error) {
	if queryString == "" {
		return nil, nil
	}

	params := strings.Split(queryString, "&")

	if len(params) == 0 {
		return nil, nil
	}

	searchParams := NewSearchParameters()

	for _, p := range params {
		parts := strings.Split(p, "=")
		if len(parts) != 2 {
			continue
		}

		key, value := qkeys.QueryKey(parts[0]), parts[1]

		if value == "" {
			continue
		}

		if !qkeys.Keys.ExistsByKey(key) {
			continue
		}

		searchParams.Add(key, value)
	}

	return searchParams, nil
}
