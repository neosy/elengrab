package qkeys

var (
	searchQueryKeys           *QueryKeys = NewQueryKeys()
	searchFilterKeys          *QueryKeys = NewQueryKeys()
	searchParameterKeys       *QueryKeys = NewQueryKeys()
	searchParameterFilterKeys *QueryKeys = NewQueryKeys()

	SearchQueryKeys           QueryKeysRegistry = searchQueryKeys
	SearchFilterKeys          QueryKeysRegistry = searchFilterKeys
	SearchParameterKeys       QueryKeysRegistry = searchParameterKeys
	SearchParameterFilterKeys QueryKeysRegistry = searchParameterFilterKeys
)

func init() {
	searchFilterKeys.Append(ChannelIDKey)

	searchParameterFilterKeys.Append(ChannelPlatformKey)

	searchParameterKeys.Append(ViewModeKey)
	searchParameterKeys.Append(LastCursorKey)
	searchParameterKeys.Append(searchParameterFilterKeys.List()...)

	searchQueryKeys.Append(SearchKey)
	searchQueryKeys.Append(SearchQueryKey)
	searchQueryKeys.Append(SearchParametersKey)
	searchQueryKeys.Append(searchFilterKeys.List()...)
	searchQueryKeys.Append(searchParameterKeys.List()...)
}
