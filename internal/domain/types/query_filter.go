package dtypes

type QueryFilter struct {
	Name  QueryFilterName
	Value any
}

type QueryFiltersByName map[QueryFilterName]QueryFilter
type QueryFiltersList []QueryFilter

func (filters *QueryFiltersList) Add(name QueryFilterName, value any) QueryFilter {
	if filters == nil {
		return QueryFilter{}
	}

	filter := QueryFilter{
		Name:  name,
		Value: value,
	}

	*filters = append(*filters, filter)

	return filter
}

func (filters QueryFiltersByName) Add(name QueryFilterName, value any) QueryFilter {
	if filters == nil {
		return QueryFilter{}
	}

	filter := QueryFilter{
		Name:  name,
		Value: value,
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
