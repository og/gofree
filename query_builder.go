package f

import (
	"fmt"
	"github.com/andreyvit/diff"
	gjson "github.com/og/x/json"
	glist "github.com/og/x/list"
	l "github.com/og/x/log"
	gmap "github.com/og/x/map"
	gtime "github.com/og/x/time"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)
type Order struct {
	Type string
	Field string
}
type Group struct {
	Field string
}
type QB struct {
	Table string
	Select []string
	Where []AND
	Offset int
	Limit int
	Order Map
	Group []string
	SoftDelete string
	Insert Map
	Update Map
	Count bool
	Debug bool
	Check []string
}


// QueryBuilder Where
type AND map[string]OP

// FindOr(Find(...), Find(...))
func Or (find  ...[]AND) (andList []AND)   {
	andList = []AND{}
	for _, v := range find {
		andList = append(andList, v[0])
	}
	return
}
//  f.And("name","nimo")
// 接收 ...interface{} 作为参数而不是 map[string]interface{} 是因为会存在这种情况
// f.And("age", f.Lt(19), "age", f.Gt(10))
func And(v ...interface{})  []AND {
	vLen := len(v)
	if vLen%2 !=0  {
		panic(errors.New("gofree: f.And(v ...inteface{}) len(v) must be even, may be you forget some."))
	}
	and := AND{}
	for i:=0;i<vLen;i++ {
		itemAny := v[i]
		var item Filter
		var isKey bool
		if i%2 == 0 { isKey = true }
		if !isKey {
			keyAny := v[i-1]
			key := keyAny.(string)
			_, has := and[key]
			_=has
			if reflect.TypeOf(itemAny).Name() != "Filter" {
				item = Eql(itemAny)
			} else {
				item = itemAny.(Filter)
			}
			if has {
				and[key] = append(and[key], item)
			} else {
				and[key] = OP{item}
			}
		}
	}
	return []AND{and}
}
func wrapField(field string) string {
	return "`" + field + "`"
}

// filter list interface maybe Filter
func (qb QB) GetSelect() (sql string, sqlValues []interface{}) {
	return qb.SQL(SQLProps{
		Statement: "SELECT",
	})
}
func (qb QB) GetUpdate() (sql string, sqlValues []interface{}) {
	return qb.SQL(SQLProps{
		Statement: "UPDATE",
	})
}
func (qb QB) GetInsert()  (sql string, sqlValues []interface{}) {
	return qb.SQL(SQLProps{
		Statement: "INSERT",
	})
}
type SQLProps struct {
	Statement string `eg:"[]string{\"SELECT\", \"UPDATE\", \"DELETE\", \"INSERT\"}"`
}
func (qb QB) SQL(props SQLProps) (sql string, sqlValues []interface{}){
	var sqlList glist.StringList
	tableName := "`" + qb.Table + "`"
	{// Statement
		switch props.Statement {
		case "SELECT":
			sqlList.Push("SELECT")
			if qb.Count {
				sqlList.Push("count(*)")
			} else {
				if len(qb.Select) == 0 {
					sqlList.Push("*")
				} else {
					sqlList.Push("`" + strings.Join(qb.Select, "`, `") + "`")
				}
			}
			sqlList.Push("FROM")
			sqlList.Push(tableName)
		case "UPDATE":
			sqlList.Push("UPDATE")
			sqlList.Push(tableName)
			sqlList.Push("SET")
			keys := gmap.UnsafeKeys(qb.Update).String()
			if len(keys) == 0 {
				panic(errors.New("gofree: update can not be a empty map"))
			}
			updateValueList := glist.StringList{}
			for _, key := range keys {
				value := qb.Update[key]
				updateValueList.Push(wrapField(key) +" = ?")
				sqlValues = append(sqlValues, value)
			}
			sqlList.Push(updateValueList.Join(", "))
		case "DELETE":
			sqlList.Push("DELETE")
		case "INSERT":
			sqlList.Push("INSERT INTO")
			sqlList.Push(tableName)
			keys := gmap.UnsafeKeys(qb.Insert).String()
			if len(keys) == 0 {
				panic(errors.New("gofree: Insert can not be a empty map"))
			}
			insertKeyList := glist.StringList{}
			insertValueList := glist.StringList{}
			for _, key := range keys {
				value := qb.Insert[key]
				insertKeyList.Push(wrapField(key))
				insertValueList.Push("?")
				sqlValues = append(sqlValues, value)
			}
			sqlList.Push("(" + insertKeyList.Join(", ") + ")")
			sqlList.Push("VALUES")
			sqlList.Push("(" + insertValueList.Join(", ") + ")")
		}
	}
	{
		// Where field operator value
		shouldWhere := len(qb.Where) != 0  || qb.SoftDelete != ""
		if props.Statement == "INSERT" {
			shouldWhere = false
		}
		if shouldWhere {
			sqlList.Push("WHERE")
			var WhereList glist.StringList
			parseWhere(qb.Where, &WhereList, &sqlValues)
			switch props.Statement {
			case "SELECT", "UPDATE":
				if qb.SoftDelete != "" {
					WhereList.Push(wrapField(qb.SoftDelete) + " IS NULL")
				}
			}
			sqlList.Push(WhereList.Join(" AND "))
			notEmptyStringSqlList := glist.StringList{}
			for _, item := range sqlList.Value {
				if item != "" {
					notEmptyStringSqlList.Push(item)
				}
			}
			sqlList = notEmptyStringSqlList
			if sqlList.Value[len(sqlList.Value)-1] == "WHERE" {
				sqlList.Pop()
			}
		}
	}
	{
		// group by
		if len(qb.Group) != 0 {
			sqlList.Push("GROUP BY")
			sqlList.Push("`" + strings.Join(qb.Group,"`, `") + "`")
		}
	}
	{
		// order by
		if len(qb.Order) != 0 {
			sqlList.Push("ORDER BY")
			orderASCList := glist.StringList{}
			orderDESCList := glist.StringList{}
			firstType := ""
			for _, field := range gmap.UnsafeKeys(qb.Order).String() {
				orderType := qb.Order[field]
				switch  orderType {
				case ASC:
					if firstType == "" {
						firstType = "ASC"
					}
					orderASCList.Push(wrapField(field))
				case DESC:
					if firstType == "" {
						firstType = "DESC"
					}
					orderDESCList.Push(wrapField(field))
				}
				orderASCList.Join(",")
			}

			orderList := glist.StringList{}
			switch firstType {
			case ASC:
				if len(orderASCList.Value) != 0 { orderList.Push(orderASCList.Join(", ") + " " + "ASC") }
				if len(orderDESCList.Value) != 0 { orderList.Push(orderDESCList.Join(", ") + " " + "DESC") }
			case DESC:
				if len(orderDESCList.Value) != 0 { orderList.Push(orderDESCList.Join(", ") + " " + "DESC") }
				if len(orderASCList.Value) != 0 { orderList.Push(orderASCList.Join(", ") + " " + "ASC") }
			}
			sqlList.Push(orderList.Join(", "))
		}
	}
	{
		// limit
		if qb.Limit != 0 && !qb.Count  {
			sqlList.Push("LIMIT ?")
			sqlValues = append(sqlValues, qb.Limit)
		}
	}
	{
		// offset
		if qb.Offset != 0 && !qb.Count {
			sqlList.Push("OFFSET ?")
			sqlValues = append(sqlValues, qb.Offset)
		}
	}
	sql = sqlList.Join(" ")
	logDebug(qb.Debug, Map{
		"sql": sql,
		"values": gjson.String(sqlValues),
	})
	if len(qb.Check) != 0 {
		matched := false
		for _, checkSQL := range qb.Check {
			if checkSQL == sql {
				matched = true
			}
		}
		if !matched {
			for _, checkSQL := range qb.Check {
				panic("sql check fail:# diff:\r\n" + diff.CharacterDiff(sql, checkSQL) + "\r\n# actual\r\n" + sql + "\r\n# expect:\r\n" + checkSQL)
			}
		}
	}

	return
}



func parseAnd (field string, op OP, whereList *glist.StringList, sqlValues *[]interface{}) {
	for _, filter := range op {
		if reflect.ValueOf(filter.Value).IsValid() && reflect.TypeOf(filter.Value).String() == "time.Time" {
			panic("gofree: can not use time.Time be sql value, mayby you should time.Format(layout), \r\n` "+ field + ":"+ filter.Value.(time.Time).Format(gtime.Second) + "`")
		}
		shouldIgnore := false
		var fieldSymbolCondition glist.StringList
		dict := filter.Dict().Kind
		switchValue := filter.Kind
		switch switchValue {
		case dict.Day:
			fieldSymbolCondition.Push(wrapField(field) + " >= ?")
			*sqlValues = append(*sqlValues, filter.TimeValue.Format(gtime.Day) + " 00:00:00")
			fieldSymbolCondition.Push("AND")
			fieldSymbolCondition.Push(wrapField(field) + " <= ?")
			*sqlValues = append(*sqlValues, filter.TimeValue.Format(gtime.Day) + " 23:59:59")
		case dict.Not:
			fieldSymbolCondition.Push(wrapField(field), "!=")
			fieldSymbolCondition.Push("?")
			*sqlValues = append(*sqlValues, filter.Value)
		case dict.IsNotNull:
			fieldSymbolCondition.Push(wrapField(field), "IS NOT NULL")
		case dict.IsNull:
			fieldSymbolCondition.Push(wrapField(field), "IS NULL")
		case dict.Custom:
			var valueList []interface{}
			anyValue := reflect.ValueOf(filter.Value)

			for i := 0; i < anyValue.Len(); i++ {
				valueList = append(valueList, anyValue.Index(i).Interface())
			}
			*sqlValues = append(*sqlValues, valueList...)
			fieldSymbolCondition.Push(wrapField(field), filter.Custom)
		case dict.CustomSQL:
			var valueList []interface{}
			anyValue := reflect.ValueOf(filter.Value)

			for i := 0; i < anyValue.Len(); i++ {
				valueList = append(valueList, anyValue.Index(i).Interface())
			}
			*sqlValues = append(*sqlValues, valueList...)
			fieldSymbolCondition.Push("(" + filter.CustomSQL + ")")
		case dict.In, dict.NotIn:
			symbol := ""
			switch switchValue {
			case dict.In:
				symbol = "IN"
			case dict.NotIn:
				symbol = "NOT IN"
			}
			fieldSymbolCondition.Push(wrapField(field), symbol)
			var valueList []interface{}
			var placeholderList glist.StringList
			anyValue := reflect.ValueOf(filter.Value)
			if anyValue.Len() == 0 {
				fieldSymbolCondition.Push("(NULL)")
			} else {
				for i := 0; i < anyValue.Len(); i++ {
					valueList = append(valueList, anyValue.Index(i).Interface())
					placeholderList.Push("?")
				}
				*sqlValues = append(*sqlValues, valueList...)
				fieldSymbolCondition.Push("(" + placeholderList.Join(", ") + ")")
			}
		case dict.Like:
			var likeValue string
			filterValueString := fmt.Sprintf("%s", filter.Value)
			switch filter.Like {
			case "start":
				likeValue = filterValueString+"%"
			case "end":
				likeValue = "%" + filterValueString
			case "have":
				likeValue = "%" + filterValueString + "%"
			}
			fieldSymbolCondition.Push(wrapField(field), "LIKE")
			fieldSymbolCondition.Push("?")
			*sqlValues = append(*sqlValues, likeValue)
		case dict.GofreeIgnore:
			shouldIgnore = true
		case dict.TimeRange:
			timeRange := filter.TimeRange.Range
			valueTime := struct {
				Start time.Time
				End time.Time
			} {}
			switch timeRange.Type {
			case timeRange.Dict().Type.Day:
				valueTime.Start = gtime.FirstHour(timeRange.Start)
				valueTime.End = gtime.LastHour(timeRange.End)
			case timeRange.Dict().Type.Month:
				valueTime.Start = gtime.FirstDay(timeRange.Start)
				valueTime.End = gtime.LastDay(timeRange.End)
			case timeRange.Dict().Type.Year:
				valueTime.Start = gtime.FirstMonth(timeRange.Start)
				valueTime.End = gtime.LastMonth(timeRange.End)
			}
			fieldSymbolCondition.Push(wrapField(field) + " >= ?")
			*sqlValues = append(*sqlValues, valueTime.Start.Format(filter.TimeRange.SQLValueLayout))
			fieldSymbolCondition.Push("AND")
			fieldSymbolCondition.Push(wrapField(field) + " <= ?")
			*sqlValues = append(*sqlValues, valueTime.End.Format(filter.TimeRange.SQLValueLayout))
		case dict.BetweenInt:
			fieldSymbolCondition.Push(wrapField(field), "BETWEEN", "?", "AND" , "?")
			valueList := []interface{}{
				filter.BetweenInt.Begin,
				filter.BetweenInt.End,
			}
			*sqlValues = append(*sqlValues, valueList...)
			break
		case dict.BetweenFloat:
			fieldSymbolCondition.Push(wrapField(field), "BETWEEN", "?", "AND" , "?")
			valueList := []interface{}{
				filter.BetweenFloat.Begin,
				filter.BetweenFloat.End,
			}
			*sqlValues = append(*sqlValues, valueList...)
			break
		case dict.BetweenString:
			fieldSymbolCondition.Push(wrapField(field), "BETWEEN", "?", "AND" , "?")
			valueList := []interface{}{
				filter.BetweenString.Begin,
				filter.BetweenString.End,
			}
			*sqlValues = append(*sqlValues, valueList...)
			break
		default:
			fieldSymbolCondition.Push(wrapField(field), filter.Symbol)
			fieldSymbolCondition.Push("?")
			*sqlValues = append(*sqlValues, filter.Value)
		}
		if !shouldIgnore {
			whereList.Push(fieldSymbolCondition.Join(" "))
		}
	}
}
func parseWhere (Where []AND, WhereList *glist.StringList, sqlValues *[]interface{}) {
	var orSqlList glist.StringList
	for _, and := range Where {
		var andList glist.StringList
		for _, field  := range gmap.UnsafeKeys(and).String() {
			op := and[field]
			parseAnd(field, op, &andList, sqlValues)
		}
		andString := andList.Join(" AND ")
		orSqlList.Push(andString)
	}
	orSqlString := orSqlList.Join(" ) OR ( ")
	if len(orSqlList.Value) > 1 {
		orSqlString = "( " + orSqlString + " )"
	}
	if orSqlString != "" {
		WhereList.Push(orSqlString)
	}
}
func logDebug(isDebug bool, data Map) {
	if !isDebug {
		return
	}
	onlyValueLogger := log.New(os.Stdout,"",log.LUTC)
	log.Print("gofree debug: ")
	for key, value := range data {
		onlyValueLogger.Print(key + ":")
		onlyValueLogger.Printf("\t%#+v",value)
	}
}
func (qb QB) BindModel(model interface{}) QB {
	if qb.Table == "" {
		tableName := reflect.ValueOf(model).MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
		if qb.Table == "" {
			panic(errors.New("tableName is empty string"))
		}
	}
	qb.SoftDelete = "deleted_at"
	return qb
}


func (qb QB) Paging(page int , perPage int) QB {
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 10
		l.V("gofree: Paging(page, perPage) alert perPage is 0 ,perPage can't be 0 . gofree will set perPage 10. but you nedd check your code.")
	}
	qb.Offset = (page - 1) * perPage
	qb.Limit = perPage
	return qb
}