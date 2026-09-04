package dtypes

type QueryOrder uint8

const (
	QueryOrderAsc QueryOrder = iota
	QueryOrderDesc
)

type QueryOrderBy struct {
	Field QueryFieldName
	Order QueryOrder
}
