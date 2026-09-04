package dbentity

type EntityInterface interface {
	// TableName returns the table name
	TableName() string

	// FieldName returns the field name from the SQL tag using a structure field name or pointer,
	// optionally prefixed with a table alias.
	//
	// Examples:
	//
	//	var entity <TableEntity>
	//	entity.FieldName("created_at", "")              // "created_at"
	//	entity.FieldName(&entity.CreatedAt, "")         // "created_at"
	//	entity.FieldName(&entity.CreatedAt, "users")    // "users.created_at"
	FieldName(field any, alias ...string) string

	// Fields returns a list of fields that will be used for queries
	FieldsAll() []string

	// FieldsAllWithAlias returns a list of fields with alias that will be used for queries
	FieldsAllWithAlias(alias string) []string

	// Fields returns a list of fields that will be used for updates
	Fields() []string
}
