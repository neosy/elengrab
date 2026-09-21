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
	list  []QueryFilter
	byKey map[qkeys.QueryKey]QueryFilter
}

func NewQueryFilters() *QueryFilters {
	return &QueryFilters{
		byKey: make(map[qkeys.QueryKey]QueryFilter),
	}
}

func (filters *QueryFilters) Append(filter QueryFilter) *QueryFilters {
	if filters == nil {
		return nil
	}

	filters.list = append(filters.list, filter)
	filters.byKey[filter.Key] = filter

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
		list:  slices.Clone(filters.list),
		byKey: maps.Clone(filters.byKey),
	}
}

func (filters *QueryFilters) Find(key qkeys.QueryKey) (QueryFilter, bool) {
	if filters == nil {
		return QueryFilter{}, false
	}

	filter, exists := filters.byKey[key]

	return filter, exists
}

func (filters *QueryFilters) GetValue(key qkeys.QueryKey) string {
	if filters == nil {
		return ""
	}

	filter, exists := filters.byKey[key]

	if !exists {
		return ""
	}

	return filter.Value
}
