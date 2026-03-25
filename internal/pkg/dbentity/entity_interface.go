package dbentity

type EntityInterface interface {
	// TableName returns the table name
	TableName() string

	// FieldName field name from sql tag by structure field name
	// Example:
	// var entity = TableEntity{}
	// entity.FieldName(&entity.SalesId)
	FieldName(field any) string

	// Fields returns a list of fields that will be used for queries
	FieldsAll() []string

	// FieldsAllWithAlias returns a list of fields with alias that will be used for queries
	FieldsAllWithAlias(alias string) []string

	// Fields returns a list of fields that will be used for updates
	Fields() []string
}
