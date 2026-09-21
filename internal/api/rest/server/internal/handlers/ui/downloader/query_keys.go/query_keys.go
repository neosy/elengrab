package qkeys

import (
	"maps"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type QueryKeys struct {
	shortKeysByKey map[string]string
	keysByShortKey map[string]string

	filterNamesByKey map[string]dtypes.QueryFilterName
	keysByFilterName map[dtypes.QueryFilterName]string
}

func newQueryKeys() *QueryKeys {
	return &QueryKeys{
		shortKeysByKey: make(map[string]string),
		keysByShortKey: make(map[string]string),

		filterNamesByKey: make(map[string]dtypes.QueryFilterName),
		keysByFilterName: make(map[dtypes.QueryFilterName]string),
	}
}

func (keys *QueryKeys) clone() *QueryKeys {
	cloneKeys := newQueryKeys()

	cloneKeys.shortKeysByKey = maps.Clone(keys.shortKeysByKey)
	cloneKeys.keysByShortKey = maps.Clone(keys.keysByShortKey)
	cloneKeys.filterNamesByKey = maps.Clone(keys.filterNamesByKey)
	cloneKeys.keysByFilterName = maps.Clone(keys.keysByFilterName)

	return cloneKeys
}

func (keys *QueryKeys) append(key QueryKey) {
	if key == "" {
		return
	}

	keys.shortKeysByKey[key.String()] = key.Short()
	keys.keysByShortKey[key.Short()] = key.String()

	if key.FilterName() != dtypes.QueryFilterNameNone {
		keys.filterNamesByKey[key.String()] = key.FilterName()
		keys.keysByFilterName[key.FilterName()] = key.String()
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

func (keys *QueryKeys) addWithName(key, shortKey string, name dtypes.QueryFilterName) QueryKey {
	keys.add(key, shortKey)

	keys.filterNamesByKey[key] = name
	keys.keysByFilterName[name] = key

	return QueryKey(key)
}

func (keys *QueryKeys) FindByKey(key QueryKey) QueryKey {
	_, exists := keys.shortKeysByKey[key.String()]
	if exists {
		return key
	}

	return QueryKey("")
}

func (keys *QueryKeys) ExistsByKey(key QueryKey) bool {
	_, exists := keys.shortKeysByKey[key.String()]
	return exists
}

func (keys *QueryKeys) FindByFilterName(name dtypes.QueryFilterName) QueryKey {
	return QueryKey(keys.keysByFilterName[name])
}

func (keys *QueryKeys) FindByShortKey(shortKey string) QueryKey {
	return QueryKey(keys.keysByShortKey[shortKey])
}

func (keys *QueryKeys) ShortKeyByKey(key QueryKey) string {
	return keys.shortKeysByKey[string(key)]
}

func (keys *QueryKeys) FilterNameByKey(key QueryKey) dtypes.QueryFilterName {
	if key == "" {
		return dtypes.QueryFilterNameNone
	}

	filterName, exists := keys.filterNamesByKey[string(key)]
	if !exists {
		return dtypes.QueryFilterNameNone
	}

	return filterName
}

func (keys *QueryKeys) FilterNameByShortKey(shortKey string) dtypes.QueryFilterName {
	if shortKey == "" {
		return dtypes.QueryFilterNameNone
	}

	key := keys.FindByShortKey(shortKey)

	return keys.FilterNameByKey(key)
}
