package f

import (
	gtime "github.com/og/x/time"
	"time"
)

type Filter struct {
	FieldWrap string
	FieldWarpArg string
	Value interface{}
	Kind string
	Symbol string
	Like string
	Custom string
	CustomSQL string
	TimeValue time.Time
	TimeRange FilterTimeRange
	BetweenInt struct{
		Begin int
		End int
	}
	BetweenFloat struct{
		Begin float64
		End float64
	}
	BetweenString struct{
		Begin string
		End string
	}
}
type FilterTimeRange struct {
	Range  gtime.Range
	SQLValueLayout string `note:"value use gtime.Day gtime.Second"`
}


func (self Filter) Dict () (dict struct {
	Kind struct{
		Day string
		Not string
		IsNotNull string
		IsNull string
		Custom string
		CustomSQL string
		In string
		NotIn string
		Like string
		GofreeIgnore string
		TimeRange string
		BetweenInt string
		BetweenFloat string
		BetweenString string
	}
}) {
	dict.Kind.Day = "Day"
	dict.Kind.Not = "NOT"
	dict.Kind.IsNotNull = "IS NOT NULL"
	dict.Kind.IsNull = "IS NULL"
	dict.Kind.Custom = "Custom"
	dict.Kind.CustomSQL = "CustomSQL"
	dict.Kind.In = "IN"
	dict.Kind.NotIn = "NOT IN"
	dict.Kind.Like = "LIKE"
	dict.Kind.GofreeIgnore = "GofreeIgnore"
	dict.Kind.TimeRange = "TimeRange"
	dict.Kind.BetweenInt = "BetweenInt"
	dict.Kind.BetweenFloat = "BetweenFloat"
	dict.Kind.BetweenString = "BetweenString"
	return
}

type FilterFunc func (v interface{}) Filter
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
func BetweenInt(begin int, end int) Filter {
	return Filter{
		BetweenInt: struct {
			Begin int
			End   int
		}{Begin: begin, End: end},
		Kind: Filter{}.Dict().Kind.BetweenInt,
	}
}
func BetweenFloat(begin float64, end float64) Filter {
	return Filter{
		BetweenFloat: struct {
			Begin float64
			End   float64
		}{Begin: begin, End: end},
		Kind: Filter{}.Dict().Kind.BetweenFloat,
	}
}
func BetweenString(begin string, end string) Filter {
	return Filter{
		BetweenString: struct {
			Begin string
			End   string
		}{Begin: begin, End: end},
		Kind: Filter{}.Dict().Kind.BetweenString,
	}
}
type nonSupportBetweenTime string
// you can use f.TimeRange, not BetweenTime
func BetweenTime(v nonSupportBetweenTime) {

}
func Like(v interface{}) Filter {
	return Filter{
		Value: v,
		Kind:Filter{}.Dict().Kind.Like,
		Like: "have",
	}
}
func LikeStart(v interface{}) Filter {
	return Filter{
		Value: v,
		Kind: Filter{}.Dict().Kind.Like,
		Like: "start",
	}
}

func LikeEnd(v interface{}) Filter {
	return Filter{
		Value: v,
		Kind: Filter{}.Dict().Kind.Like,
		Like: "end",
	}
}

func CustomSQL(sql string, values ...interface{}) Filter {
	return Filter{
		Value: values,
		Kind: Filter{}.Dict().Kind.CustomSQL,
		CustomSQL: sql,
	}
}
func Custom (template string, v ...interface{}) Filter {
	return Filter{
		Value: v,
		Kind: Filter{}.Dict().Kind.Custom,
		Custom: template,
	}
}
func In (v interface{}) Filter {
	return Filter{
		Value: v,
		Kind: Filter{}.Dict().Kind.In,
	}
}

func NotIn (v interface{}) Filter {
	return Filter{
		Value: v,
		Kind: Filter{}.Dict().Kind.NotIn,
	}
}
func IsNull () Filter {
	return Filter{
		Kind: Filter{}.Dict().Kind.IsNull,
	}
}

func IsNotNull () Filter {
	return Filter{
		Kind: Filter{}.Dict().Kind.IsNotNull,
	}
}
func Day(v time.Time) Filter{
	return Filter{
		Kind: Filter{}.Dict().Kind.Day,
		TimeValue: v,
	}
}



func IgnoreFilter () Filter {
	return Filter{
		Kind: Filter{}.Dict().Kind.GofreeIgnore,
	}
}
type orderType uint8
const (
	DESC orderType = iota
	ASC
)

// 在查询中有一种常见的场景，当某个请求参数为空时不增加 where。
// 比如用户搜索姓名, ?name=nimo 时SQL是 WHERE name = ? 。
// 如果 ?name= （空字符串）则 sql 没有 name = ?
// gofree 称这种 where 条件为 ignore empty
/*
	使用场景:
	f.QB{
		Where: f.And(
			"name": f.IgnoreEmpty(f.Eql, query.Name)
		),
	}
*/
func IgnoreEmpty(filterFunc FilterFunc, query string) Filter {
	return IgnorePattern(filterFunc, query, "")
}


// 基于 IgnoreEmpty 的场景下，有些请求并不一定会是空，而是 ?status=all 来表示搜索全部
// ?status=done 表示搜索已完成的数据 ,此时使用 IgnorePattern(f.Eql, query.Status, "all")
/*
	使用场景:
	f.QB{
		Where: f.And(
			"status": f.IgnorePattern( f.Eql, query.Status, "all")
		),
	}
*/
func IgnorePattern(filterFunc FilterFunc, query string, pattern string) Filter {
	if query == pattern {
		return IgnoreFilter()
	} else {
		return filterFunc(query)
	}
}
// 在 IgnoreEmpty 和 IgnorePattern 的场景下WHERE 语句都是 field = ?
// 有些场景下可能需要 where field in ? 或者没有 field in ?
// 此时使用 Ignore 完全自定义控制
/*
	使用场景:
	f.QB{
		Where: f.And(
			"id": f.Ignore(len(query.idList) == 0, f.In(query.idList))
		),
	}
*/
func Ignore(filter Filter, ignoreCondition bool)  Filter {
	if ignoreCondition {
		return IgnoreFilter()
	} else {
		return filter
	}
}

func TimeRange(data gtime.Range, sqlValueLayout string) Filter {
	return Filter{
		Kind: Filter{}.Dict().Kind.TimeRange,
		TimeRange: FilterTimeRange{Range: data, SQLValueLayout: sqlValueLayout},
	}
}


