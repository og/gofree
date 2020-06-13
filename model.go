package f

type Model interface {
	TableName () string
	BeforeCreate()
}

type RelationModel interface {
	Relation () (tableName string, join []Join)
}
































