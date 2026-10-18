package pages

import (
	"slices"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

var (
	visibilityList []VisibilityOption
)

func init() {
	visibilityList = newVisibilityList()
}

func VisibilityList() []VisibilityOption {
	return slices.Clone(visibilityList)
}

func newVisibilityList() []VisibilityOption {
	var list = make([]VisibilityOption, len(dtypes.MediaVisibilityList()))
	for i, v := range dtypes.MediaVisibilityList() {
		list[i] = VisibilityOption{
			Value: v.String(),
			Label: v.Label(),
		}
	}
	return list
}
