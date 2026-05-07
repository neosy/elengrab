package tablenames

const (
	Links      = "links"
	LinkClicks = "link_clicks"
)

var tableNames = []string{
	Links,
	LinkClicks,
}

func TableNames() []string {
	return tableNames
}
