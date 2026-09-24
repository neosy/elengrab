package qkeys

import (
	"maps"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type QueryKeysRegistry interface {
	FindByKey(key QueryKey) QueryKey
	ExistsByKey(key QueryKey) bool

	FindByStringKey(key string) QueryKey
	ExistsByStringKey(key string) bool

	FindByShortKey(shortKey string) QueryKey
	FindShortKeyByKey(key QueryKey) string

	FindByFilterName(name dtypes.QueryFilterName) QueryKey
	FilterNameByKey(key QueryKey) dtypes.QueryFilterName
	FilterNameByShortKey(shortKey string) dtypes.QueryFilterName

	Len() int
	List() []QueryKey

	// Intersect returns a new QueryKeys containing only the keys from inKeys
	// that exist in the current keys collection. Returns nil if either list is nil
	Intersect(inKeys QueryKeysRegistry) *QueryKeys
}

type QueryKeys struct {
	shortKeysByKey map[QueryKey]string
	keysByShortKey map[string]QueryKey

	filterNamesByKey map[string]dtypes.QueryFilterName
	keysByFilterName map[dtypes.QueryFilterName]string
}

func NewQueryKeys() *QueryKeys {
	return &QueryKeys{
		shortKeysByKey: make(map[QueryKey]string),
		keysByShortKey: make(map[string]QueryKey),

		filterNamesByKey: make(map[string]dtypes.QueryFilterName),
		keysByFilterName: make(map[dtypes.QueryFilterName]string),
	}
}

func (keys *QueryKeys) Clone() *QueryKeys {
	cloneKeys := NewQueryKeys()

	cloneKeys.shortKeysByKey = maps.Clone(keys.shortKeysByKey)
	cloneKeys.keysByShortKey = maps.Clone(keys.keysByShortKey)
	cloneKeys.filterNamesByKey = maps.Clone(keys.filterNamesByKey)
	cloneKeys.keysByFilterName = maps.Clone(keys.keysByFilterName)

	return cloneKeys
}

func (keys *QueryKeys) Append(qKeys ...QueryKey) {
	if len(qKeys) == 0 {
		return
	}

	add := func(key QueryKey) {
		keys.shortKeysByKey[key] = key.Short()
		keys.keysByShortKey[key.Short()] = key

		if key.FilterName() != dtypes.QueryFilterNameNone {
			keys.filterNamesByKey[key.String()] = key.FilterName()
			keys.keysByFilterName[key.FilterName()] = key.String()
		}
	}

	for _, key := range qKeys {
		if key == "" {
			continue
		}

		add(key)
	}
}

func (keys *QueryKeys) add(key, shortKey string) QueryKey {
	_, exists := keys.shortKeysByKey[QueryKey(key)]
	if exists {
		panic("query key already exists: " + key)
	}

	_, exists = keys.keysByShortKey[shortKey]
	if exists {
		panic("query short key already exists: " + shortKey)
	}

	keys.shortKeysByKey[QueryKey(key)] = shortKey
	keys.keysByShortKey[shortKey] = QueryKey(key)

	return QueryKey(key)
}

func (keys *QueryKeys) addWithName(key, shortKey string, name dtypes.QueryFilterName) QueryKey {
	keys.add(key, shortKey)

	keys.filterNamesByKey[key] = name
	keys.keysByFilterName[name] = key

	return QueryKey(key)
}

func (keys *QueryKeys) FindByKey(key QueryKey) QueryKey {
	_, exists := keys.shortKeysByKey[key]
	if exists {
		return key
	}

	return QueryKey("")
}

func (keys *QueryKeys) FindByStringKey(key string) QueryKey {
	return keys.FindByKey(QueryKey(key))
}

func (keys *QueryKeys) ExistsByKey(key QueryKey) bool {
	_, exists := keys.shortKeysByKey[key]
	return exists
}

func (keys *QueryKeys) ExistsByStringKey(key string) bool {
	return keys.ExistsByKey(QueryKey(key))
}

func (keys *QueryKeys) FindByFilterName(name dtypes.QueryFilterName) QueryKey {
	return QueryKey(keys.keysByFilterName[name])
}

func (keys *QueryKeys) FindByShortKey(shortKey string) QueryKey {
	return QueryKey(keys.keysByShortKey[shortKey])
}

func (keys *QueryKeys) FindShortKeyByKey(key QueryKey) string {
	return keys.shortKeysByKey[key]
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

func (keys *QueryKeys) Len() int {
	return len(keys.shortKeysByKey)
}

func (keys *QueryKeys) List() []QueryKey {
	var qKeys []QueryKey

	for key := range keys.shortKeysByKey {
		qKeys = append(qKeys, key)
	}

	return qKeys
}

// Intersect returns a new QueryKeys containing only the keys from inKeys
// that exist in the current keys collection. Returns nil if either list is nil
func (keys *QueryKeys) Intersect(inKeys QueryKeysRegistry) *QueryKeys {
	if keys == nil || inKeys == nil {
		return nil
	}

	outKeys := NewQueryKeys()

	for _, key := range inKeys.List() {
		if keys.ExistsByKey(key) {
			outKeys.Append(key)
		}
	}

	if outKeys.Len() == 0 {
		return nil
	}

	return outKeys
}
