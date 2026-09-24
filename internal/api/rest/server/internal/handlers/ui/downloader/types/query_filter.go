package types

import (
	"maps"
	"slices"

	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
)

type QueryFilter struct {
	Key   qkeys.QueryKey
	Value string
}

type QueryFilters struct {
	list        []QueryFilter
	valuesByKey map[qkeys.QueryKey]string
}

func NewQueryFilters() *QueryFilters {
	return &QueryFilters{
		valuesByKey: make(map[qkeys.QueryKey]string),
	}
}

func (filters *QueryFilters) Append(filter QueryFilter) *QueryFilters {
	if filters == nil {
		return nil
	}

	filters.list = append(filters.list, filter)
	filters.valuesByKey[filter.Key] = filter.Value

	return filters
}

func (filters *QueryFilters) Add(key qkeys.QueryKey, value string) QueryFilter {
	if filters == nil {
		return QueryFilter{}
	}

	filter := QueryFilter{
		Key:   key,
		Value: value,
	}

	filters.Append(filter)

	return filter
}

func (filters *QueryFilters) List() []QueryFilter {
	if filters == nil {
		return nil
	}

	return slices.Clone(filters.list)
}

func (filters *QueryFilters) Len() int {
	if filters == nil {
		return 0
	}

	return len(filters.list)
}

func (filters *QueryFilters) Clone() *QueryFilters {
	if filters == nil {
		return nil
	}

	return &QueryFilters{
		list:        slices.Clone(filters.list),
		valuesByKey: maps.Clone(filters.valuesByKey),
	}
}

func (filters *QueryFilters) Find(key qkeys.QueryKey) (QueryFilter, bool) {
	if filters == nil {
		return QueryFilter{}, false
	}

	value, exists := filters.valuesByKey[key]

	return QueryFilter{key, value}, exists
}

func (filters *QueryFilters) Exists(key qkeys.QueryKey) bool {
	if filters == nil {
		return false
	}

	_, exists := filters.valuesByKey[key]

	return exists
}

func (filters *QueryFilters) GetValue(key qkeys.QueryKey) string {
	if filters == nil {
		return ""
	}

	value, exists := filters.valuesByKey[key]

	if !exists {
		return ""
	}

	return value
}

func (filters *QueryFilters) QueryKeys() *qkeys.QueryKeys {
	if filters == nil || filters.Len() == 0 {
		return nil
	}

	keys := qkeys.NewQueryKeys()

	for _, filter := range filters.list {
		keys.Append(filter.Key)
	}

	return keys
}

// FilterByKeys filters the current collection and returns a new QueryFilters
// containing only the filters whose keys exist in the provided keys registry.
// Returns nil if either the current collection or the keys registry is nil.
func (filters *QueryFilters) FilterByKeys(keys qkeys.QueryKeysRegistry) *QueryFilters {
	if filters == nil || keys == nil {
		return nil
	}

	outFilters := NewQueryFilters()

	for _, filter := range filters.List() {
		if keys.ExistsByKey(filter.Key) {
			outFilters.Append(filter)
		}
	}

	return outFilters
}
