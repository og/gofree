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
func ignoreFilter () Filter {
	return Filter{
		Symbol: "GOFREE_IGNORE",
	}
}
const DESC = "DESC"
const ASC = "ASC"

// 在查询中有一种常见的场景，当某个请求参数为空时不增加 where。
// 比如用户搜索姓名, ?name=nimo 时SQL是 WHERE name = ? 。
// 如果 ?name= （空字符串）则 sql 没有 name = ?
// gofree 称这种 where 条件为 ignore empty
func IgnoreEmpty(query string) Filter {
	return IgnorePattern(query, "")
}

// 基于 IgnoreEmpty 的场景下，有些请求并不一定会是空，而是 ?status=all 来表示搜索全部
// ?status=done 表示搜索已完成的数据 ,此时使用 IgnorePattern(query, "all")
func IgnorePattern(query string, pattern string) Filter {
	if query == pattern {
		return ignoreFilter()
	} else {
		return Eql(query)
	}
}


// 在 IgnoreEmpty 和 IgnorePattern 的场景下WHERE 语句都是 field = ?
// 有些场景下可能需要 where field in ? 或者没有 field in ?
// 此时使用 Ignore 完全自定义控制
func Ignore(ignoreCondition bool, filter Filter)  Filter {
	if ignoreCondition {
		return ignoreFilter()
	} else {
		return filter
	}
}
