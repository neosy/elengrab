package types

type SearchValues struct {
	QueryText  string
	Filters    *QueryFilters
	Parameters SearchParameterValues
}

func NewSearchValues() SearchValues {
	return SearchValues{
		Parameters: NewSearchParameterValues(),
	}
}

func (v *SearchValues) IsZero() bool {
	if v == nil {
		return true
	}

	return v.QueryText == "" &&
		v.Filters.Len() == 0 &&
		v.Parameters.IsZero()
}
