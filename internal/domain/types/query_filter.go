package dtypes

import "github.com/neosy/elengrab/internal/pkg/dbutils"

type QueryFilter struct {
	Name      QueryFilterName
	condition dbutils.FilterConditioner
}

type QueryFiltersByName map[QueryFilterName]QueryFilter
type QueryFiltersList []QueryFilter

func NewQueryFilter(name QueryFilterName, condition dbutils.FilterConditioner) QueryFilter {
	if condition == nil {
		condition = dbutils.NewFilterConditionEq[any](nil)
	}

	return QueryFilter{
		Name:      name,
		condition: condition,
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

func (filters *QueryFiltersList) Append(filter QueryFilter) QueryFilter {
	*filters = append(*filters, filter)
	return filter
}

func (filters *QueryFiltersList) Add(name QueryFilterName, value any) QueryFilter {
	if filters == nil {
		return QueryFilter{}
	}

	filter := QueryFilter{
		Name:      name,
		condition: dbutils.NewFilterCondition(value, dbutils.FilterOperatorEq),
	}

	*filters = append(*filters, filter)

	return filter
}

func (filters QueryFiltersByName) Add(name QueryFilterName, value any) QueryFilter {
	if filters == nil {
		return QueryFilter{}
	}

	filter := QueryFilter{
		Name:      name,
		condition: dbutils.NewFilterCondition(value, dbutils.FilterOperatorEq),
	}

	filters[name] = filter

	return filter
}

func (filters QueryFiltersByName) List() QueryFiltersList {
	list := make(QueryFiltersList, len(filters))
	for _, filter := range filters {
		list = append(list, filter)
	}
	return list
}

func (filters QueryFiltersList) FiltersByName() QueryFiltersByName {
	filtersByName := make(QueryFiltersByName, len(filters))
	for _, filter := range filters {
		filtersByName[filter.Name] = filter
	}
	return filtersByName
}
