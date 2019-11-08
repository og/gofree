package f

import (
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
func Day(v string) Filter{
	return Filter{
		FieldWrap: "DATE",
		Symbol: "day",
		Value: v,
	}
}
func Month(v string) Filter {
	return Filter{
		FieldWrap: "DATE_FORMAT",
		FieldWarpArg: "%Y%m",
		Symbol: "month",
		Value: v,
	}
}
func Year(v string) Filter {
	return Filter{
		FieldWrap: "DATE_FORMAT",
		FieldWarpArg: "%Y",
		Symbol: "year",
		Value: v,
	}
}
const DESC = "DESC"
const ASC = "ASC"
