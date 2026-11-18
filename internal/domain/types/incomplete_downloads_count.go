package dtypes

type IncompleteDownloadsCount struct {
	UserInProgressCount int `json:"userInProgressCount"`
	UserFailedCount     int `json:"userFailedCount"`
	UserCount           int `json:"userCount"`

	TotalInProgressCount int `json:"totalInProgressCount"`
	TotalFailedCount     int `json:"totalFailedCount"`
	TotalCount           int `json:"totalCount"`
}

func (c *IncompleteDownloadsCount) ClearTotal() {
	c.TotalInProgressCount = 0
	c.TotalFailedCount = 0
	c.TotalCount = 0
}

func (c *IncompleteDownloadsCount) Add(other IncompleteDownloadsCount) {
	c.UserInProgressCount += other.UserInProgressCount
	c.UserFailedCount += other.UserFailedCount
	c.UserCount += other.UserCount

	c.TotalInProgressCount += other.TotalInProgressCount
	c.TotalFailedCount += other.TotalFailedCount
	c.TotalCount += other.TotalCount
}
