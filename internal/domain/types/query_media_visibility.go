package dtypes

type QueryMediaVisibility uint8

const (
	QueryMediaVisibilityAll QueryMediaVisibility = iota
	QueryMediaVisibilityPublic
	QueryMediaVisibilityAuthenticated
)
