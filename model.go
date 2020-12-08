package f

type Model interface {
	TableName () string
	BeforeCreate()
}
type Column string
func NewColumn(column string) Column {
	return Column(column)
}
func (c Column) String() string {
	return string(c)
}
type AutoIncrement uint

type RelationModel interface {
	TableName () string
	RelationJoin () []Join
}



























