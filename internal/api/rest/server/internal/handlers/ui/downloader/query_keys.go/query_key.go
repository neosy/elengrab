package qkeys

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type QueryKey string

func (k QueryKey) String() string {
	return string(k)
}

func (k QueryKey) Short() string {
	return Keys.ShortKeyByKey(k)
}

func (k QueryKey) FilterName() dtypes.QueryFilterName {
	return Keys.FilterNameByKey(k)
}

func (k QueryKey) IsZero() bool {
	return k == ""
}
