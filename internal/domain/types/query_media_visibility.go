package dtypes

type QueryMediaVisibility uint8

const (
	// QueryMediaVisibilityAll includes media of all visibility levels.
	QueryMediaVisibilityAll QueryMediaVisibility = iota

	// QueryMediaVisibilityPublic includes only publicly visible media.
	QueryMediaVisibilityPublic

	// QueryMediaVisibilityAuthenticated includes media visible to authenticated users.
	QueryMediaVisibilityAuthenticated
)
