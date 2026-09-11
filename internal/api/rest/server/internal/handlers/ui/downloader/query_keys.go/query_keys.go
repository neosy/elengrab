package qkeys

type QueryKeys struct {
	shortKeysByKey map[string]string
	keysByShortKey map[string]string
}

func NewQueryKeys() *QueryKeys {
	return &QueryKeys{
		shortKeysByKey: make(map[string]string),
		keysByShortKey: make(map[string]string),
	}
}

func (keys *QueryKeys) add(key, shortKey string) QueryKey {
	_, exists := keys.shortKeysByKey[key]
	if exists {
		panic("query key already exists: " + key)
	}

	_, exists = keys.keysByShortKey[shortKey]
	if exists {
		panic("query short key already exists: " + shortKey)
	}

	keys.shortKeysByKey[key] = shortKey
	keys.keysByShortKey[shortKey] = key

	return QueryKey(key)
}

func (keys *QueryKeys) ShortKeyByKey(key string) string {
	return keys.shortKeysByKey[key]
}

func (keys *QueryKeys) KeyByShortKey(shortKey string) QueryKey {
	return QueryKey(keys.keysByShortKey[shortKey])
}
