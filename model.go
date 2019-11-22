package f

type Model interface {
	TableName () string
	BeforeCreate()
}
