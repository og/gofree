package f

import (
	"time"
)

type OP []Filter
type Filter struct {
	FieldWrap string
	FieldWarpArg string
	Value interface{}
	Symbol string
	Like string
	Custom string
	CustomSQL string
	TimeValue time.Time
}
func Eql(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "=",
	}
}
func NotEql(v interface{}) Filter{
	return Filter{
		Value: v,
		Symbol: "!=",
	}
}
func Lt(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "<",
	}
}
func LtEql(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "<=",
	}
}
func Gt(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: ">",
	}
}
func GtEql(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: ">=",
	}
}
func Like(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "LIKE",
		Like: "have",
	}
}
func LikeStart(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "LIKE",
		Like: "start",
	}
}

func LikeEnd(v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "LIKE",
		Like: "end",
	}
}

func CustomSQL(sql string, values ...interface{}) Filter {
	return Filter{
		Value: values,
		Symbol: "CustomSQL",
		CustomSQL: sql,
	}
}
func Custom (template string, v ...interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "custom",
		Custom: template,
	}
}
func In (v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "IN",
	}
}

func NotIn (v interface{}) Filter {
	return Filter{
		Value: v,
		Symbol: "NOT IN",
	}
}
func IsNull () Filter {
	return Filter{
		Symbol: "IS NULL",
	}
}

func IsNotNull () Filter {
	return Filter{
		Symbol: "IS NOT NULL",
	}
}
func Day(v time.Time) Filter{
	return Filter{
		Symbol: "day",
		TimeValue: v,
	}
}
const DESC = "DESC"
const ASC = "ASC"
