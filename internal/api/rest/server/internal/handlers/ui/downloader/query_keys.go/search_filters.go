package qkeys

var (
	searchQueryKeys *QueryKeys
)

func init() {
	searchQueryKeys = newQueryKeys()

	searchQueryKeys.append(ChannelIDKey)
}

func SearchKeys() *QueryKeys {
	return searchQueryKeys.clone()
}
