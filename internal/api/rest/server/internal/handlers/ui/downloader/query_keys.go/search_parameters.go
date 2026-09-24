package qkeys

var (
	searchFilterKeys    *QueryKeys = NewQueryKeys()
	searchParameterKeys *QueryKeys = NewQueryKeys()
	searchQueryKeys     *QueryKeys = NewQueryKeys()

	SearchFilterKeys    QueryKeysRegistry = searchFilterKeys
	SearchParameterKeys QueryKeysRegistry = searchParameterKeys
	SearchQueryKeys     QueryKeysRegistry = searchQueryKeys
)

func init() {
	searchFilterKeys.Append(ChannelIDKey)

	searchParameterKeys.Append(searchFilterKeys.List()...)
	searchParameterKeys.Append(ViewModeKey)
	searchParameterKeys.Append(LastCursorKey)

	searchQueryKeys.Append(searchParameterKeys.List()...)
	searchQueryKeys.Append(SearchKey)
	searchQueryKeys.Append(SearchQueryKey)
	searchQueryKeys.Append(SearchParametersKey)
}
