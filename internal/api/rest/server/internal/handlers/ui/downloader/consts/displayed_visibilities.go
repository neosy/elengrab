package consts

import dtypes "github.com/neosy/elengrab/internal/domain/types"

var DisplayedVisibilities = map[dtypes.MediaVisibility]struct{}{
	dtypes.MediaVisibilityPrivate:       {},
	dtypes.MediaVisibilityAuthenticated: {},
	dtypes.MediaVisibilityPublic:        {},
}
