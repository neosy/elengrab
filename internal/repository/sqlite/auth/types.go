package auth

import "maps"

type filtersByName map[string]any

func (filters *filtersByName) copy() filtersByName {
	if filters == nil {
		return nil
	}

	newFilters := make(filtersByName)
	maps.Copy(newFilters, *filters)

	return newFilters
}
