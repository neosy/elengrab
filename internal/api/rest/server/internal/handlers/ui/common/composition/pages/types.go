package pages

import dtypes "github.com/neosy/elengrab/internal/domain/types"

type ViewModeTab struct {
	Label  string
	Value  string
	Active bool
}

type VisibilityOption struct {
	Value string
	Label string
}

func BuildViewModeTabs(activeMode dtypes.QueryMediaViewMode) []ViewModeTab {
	var tabs []ViewModeTab

	for _, mode := range dtypes.ViewModeList() {
		tab := ViewModeTab{
			Label:  mode.Label(),
			Value:  mode.String(),
			Active: mode == activeMode,
		}

		tabs = append(tabs, tab)
	}

	return tabs
}
