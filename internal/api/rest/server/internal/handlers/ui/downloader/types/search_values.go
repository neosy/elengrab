package types

type SearchValues struct {
	QueryText  string
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
		v.Parameters.IsZero()
}
