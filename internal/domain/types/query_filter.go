package dtypes

import (
	"maps"
	"slices"
	"sync"

	"github.com/neosy/elengrab/internal/pkg/dbutils"
)

type QueryFilter struct {
	Name      QueryFilterName
	condition dbutils.FilterConditioner
}

type QueryFilters struct {
	mu sync.RWMutex

	list   []QueryFilter
	byName map[QueryFilterName]QueryFilter
}

func NewQueryFilter(name QueryFilterName, condition dbutils.FilterConditioner) QueryFilter {
	if condition == nil {
		condition = dbutils.NewFilterConditionEq[any](nil)
	}

	return QueryFilter{
		Name:      name,
		condition: condition,
	}
}

func NewQueryFilters() *QueryFilters {
	return &QueryFilters{
		byName: make(map[QueryFilterName]QueryFilter),
	}
}

func (f QueryFilter) Condition() dbutils.FilterConditioner {
	return f.condition
}

func (f QueryFilter) Value() any {
	if f.condition == nil {
		return nil
	}
	return f.condition.Value()
}

func (filters *QueryFilters) Append(filter QueryFilter) *QueryFilters {
	filters.mu.Lock()
	defer filters.mu.Unlock()

	filters.list = append(filters.list, filter)
	filters.byName[filter.Name] = filter

	return filters
}

func (filters *QueryFilters) Add(name QueryFilterName, value any) QueryFilter {
	filter := QueryFilter{
		Name:      name,
		condition: dbutils.NewFilterCondition(value, dbutils.FilterOperatorEq),
	}

	filters.Append(filter)

	return filter
}

func (filters *QueryFilters) List() []QueryFilter {
	if filters == nil {
		return nil
	}

	filters.mu.RLock()
	defer filters.mu.RUnlock()

	return slices.Clone(filters.list)
}

func (filters *QueryFilters) Len() int {
	if filters == nil {
		return 0
	}

	filters.mu.RLock()
	defer filters.mu.RUnlock()

	return len(filters.list)
}

func (filters *QueryFilters) Clone() *QueryFilters {
	if filters == nil {
		return nil
	}

	filters.mu.RLock()
	defer filters.mu.RUnlock()

	return &QueryFilters{
		list:   slices.Clone(filters.list),
		byName: maps.Clone(filters.byName),
	}
}

func (filters *QueryFilters) Find(name QueryFilterName) (QueryFilter, bool) {
	filter, exists := filters.byName[name]

	return filter, exists
}

func (filters *QueryFilters) GetValue(name QueryFilterName) any {
	filter, exists := filters.byName[name]

	if !exists {
		return nil
	}

	return filter.Value()
}

func (filters *QueryFilters) GetStringValue(name QueryFilterName) string {
	value := filters.GetValue(name)
	if value == nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return ""
	}

	return strValue
}
