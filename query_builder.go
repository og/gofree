package f

import (
	"errors"
	"fmt"
	"github.com/andreyvit/diff"
	gjson "github.com/og/x/json"
	gmap "github.com/og/x/map"
	gtime "github.com/og/x/time"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)
type Order struct {
	// Field Type 的顺序永远不能换
	Column Column
	Type orderType
}
type Group struct {
	Column Column
}
type QB struct {
	Table string
	Select []Column
	Where []WhereAnd
	Offset int
	Limit int
	Order []Order
	Group []Column
	SoftDelete Column
	Insert Data
	Update Data
	Count bool
	Debug bool
	Check []string
}
type WhereAnd map[Column][]Filter
type WhereAndList []WhereAnd
func filterListRemove(list []Filter, handle func(index int, filter Filter) (remove bool)) (newFilter []Filter) {
	for filterIndex, filterItem := range list {
		if !handle(filterIndex, filterItem) {
			newFilter = append(newFilter, filterItem)
		}
	}
	return
}
func whereAndListRemove(list []WhereAnd, handle func(index int, and WhereAnd) (remove bool)) (newWhere []WhereAnd) {
	for andIndex, andItem := range list {
		if !handle(andIndex, andItem) {
			newWhere = append(newWhere, andItem)
		}
	}
	return
}
func whereAndListRemoveByIndex(list []WhereAnd, removeIndex int) (newWhere []WhereAnd) {
	return whereAndListRemove(list, func(index int, and WhereAnd) (remove bool) {
		return index == removeIndex
	})
}
func (where WhereAndList) And(column Column, value interface{}) WhereAndList {
	var filterValue Filter
	switch value.(type) {
	case Filter:
		filterValue = value.(Filter)
	default:
		filterValue = Equal(value)
	}
	if len(where) == 0 {
		where = append(where, map[Column][]Filter{})
	}
	where[0][column] = append(where[0][column], filterValue)
	return where
}
func And(column Column, value interface{}) WhereAndList {
	return WhereAndList{}.And(column, value)
}
func wrapField(field Column) string {
	return "`" + string(field) + "`"
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
	Statement statement `eg:"[]string{\"SELECT\", \"UPDATE\", \"DELETE\", \"INSERT\"}"`
}
func (SQLProps) Dict() (dict struct {
	Statement struct {
		Select statement
		Update statement
		Delete statement
		INSERT statement
	}
}) {
	dict.Statement.Select = "SELECT"
	dict.Statement.Update = "UPDATE"
	dict.Statement.Delete = "DELETE"
	dict.Statement.INSERT = "INSERT"
	return
}
func (v SQLProps) SwitchStatement(
	Insert func(_Insert int),
	Select func(_Select bool),
	Update func(_Update string),
	Delete func(_Delete []int),
	) {
	dict := v.Dict().Statement
	switch v.Statement {
	default:
		panic("gofree: f.SQLProps{}.SwitchStatement Statement switch error")
	case dict.INSERT:
		Insert(1)
	case dict.Select:
		Select(true)
	case dict.Update:
		Update("")
	case dict.Delete:
		Delete([]int{1})
	}
}
type statement string
func (qb QB) SQL(props SQLProps) (sql string, sqlValues []interface{}){
	var sqlList stringQueue
	tableName := "`" + qb.Table + "`"
	{// Statement
		props.SwitchStatement(
			func(_Insert int) {
				sqlList.Push("INSERT INTO")
				sqlList.Push(tableName)
				var keys []Column
				for _, key := range gmap.UnsafeKeys(qb.Insert).String() {
					keys = append(keys, Column(key))
				}
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
		},
		func(_Select bool) {
			sqlList.Push("SELECT")
			if qb.Count {
				sqlList.Push("count(*)")
			} else {
				if len(qb.Select) == 0 {
					sqlList.Push("*")
				} else {
					sqlList.Push("`" + strings.Join(ColumnsToStrings(qb.Select), "`, `") + "`")
				}
			}
			sqlList.Push("FROM")
			sqlList.Push(tableName)
		}, func(_Update string) {
			sqlList.Push("UPDATE")
			sqlList.Push(tableName)
			sqlList.Push("SET")
			var keys []Column
			for _, key := range gmap.UnsafeKeys(qb.Update).String() {
				keys = append(keys, Column(key))
			}
			if len(keys) == 0 {
				panic(errors.New("gofree: update can not be a empty map"))
			}
			updateValueList := stringQueue{}
			for _, key := range keys {
				value := qb.Update[key]
				updateValueList.Push(wrapField(key) +" = ?")
				sqlValues = append(sqlValues, value)
			}
			sqlList.Push(updateValueList.Join(", "))
		}, func(_Delete []int) {
				sqlList.Push("DELETE")
		})
	}
	{
		// Where field operator value
		// remove ignore
		for _, whereAnd := range qb.Where {
			for column, filterList := range whereAnd {
				for _, filter := range filterList {
					if filter.Kind == filter.Dict().Kind.GofreeIgnore {
						if len(filterList) == 1 {
							delete(whereAnd, column)
						}
					}
				}
			}
		}
		for whereAndIndex, whereAnd := range qb.Where {
			if len(whereAnd) == 0 {
				qb.Where = whereAndListRemoveByIndex(qb.Where, whereAndIndex)
			}
		}
		shouldWhere := len(qb.Where) != 0  || qb.SoftDelete != ""
		if props.Statement == "INSERT" {
			shouldWhere = false
		}
		if shouldWhere {
			sqlList.Push("WHERE")
			var WhereList stringQueue
			parseWhere(qb.Where, &WhereList, &sqlValues)
			switch props.Statement {
			case "SELECT", "UPDATE":
				if qb.SoftDelete != "" {
					WhereList.Push(wrapField(Column(qb.SoftDelete)) + " IS NULL")
				}
			}
			sqlList.Push(WhereList.Join(" AND "))
			notEmptyStringSqlList := stringQueue{}
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
			var gourpString []string
			for _, group := range qb.Group {
				gourpString = append(gourpString, string(group))
			}
			sqlList.Push("`" + strings.Join(gourpString,"`, `") + "`")
		}
	}
	{
		// order by
		if len(qb.Order) != 0 {
			sqlList.Push("ORDER BY")
			orderList := stringQueue{}
			for _, orderItem := range qb.Order {
				orderType := orderItem.Type
				switch  orderType {
				case ASC:
					orderList.Push(wrapField(orderItem.Column) +" ASC")
				case DESC:
					orderList.Push(wrapField(orderItem.Column)+" DESC")
				}
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
	logDebug(qb.Debug, Data{
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



func parseAnd (field string, op OP, whereList *stringQueue, sqlValues *[]interface{}) {
	for _, filter := range op {
		if reflect.ValueOf(filter.Value).IsValid() && reflect.TypeOf(filter.Value).String() == "time.Time" {
			panic("gofree: can not use time.Time be sql value, mayby you should time.Format(layout), \r\n` "+ field + ":"+ filter.Value.(time.Time).Format(gtime.LayoutTime) + "`")
		}
		shouldIgnore := false
		var fieldSymbolCondition stringQueue
		dict := filter.Dict().Kind
		switchValue := filter.Kind
		switch switchValue {
		case dict.Day:
			fieldSymbolCondition.Push(wrapField(field) + " >= ?")
			*sqlValues = append(*sqlValues, filter.TimeValue.Format(gtime.LayoutDate) + " 00:00:00")
			fieldSymbolCondition.Push("AND")
			fieldSymbolCondition.Push(wrapField(field) + " <= ?")
			*sqlValues = append(*sqlValues, filter.TimeValue.Format(gtime.LayoutDate) + " 23:59:59")
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
			var placeholderList stringQueue
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
			timeRange.Type.Switch(func(_year int) {
				valueTime.Start = gtime.FirstMonth(timeRange.Start)
				valueTime.End = gtime.LastMonth(timeRange.End)
			}, func(_month bool) {
				valueTime.Start = gtime.FirstDay(timeRange.Start)
				valueTime.End = gtime.LastDay(timeRange.End)
			}, func(_day string) {
				valueTime.Start = gtime.FirstHour(timeRange.Start)
				valueTime.End = gtime.LastHour(timeRange.End)
			})
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
			if filter.Symbol == "" {
				panic(errors.New("f.Filter is empty struct"))
			}
			fieldSymbolCondition.Push(wrapField(field), filter.Symbol)
			fieldSymbolCondition.Push("?")
			*sqlValues = append(*sqlValues, filter.Value)
		}
		if !shouldIgnore {
			whereList.Push(fieldSymbolCondition.Join(" "))
		}
	}
}
func parseWhere (Where []AND, WhereList *stringQueue, sqlValues *[]interface{}) {
	var orSqlList stringQueue
	for _, and := range Where {
		var andList stringQueue
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
func logDebug(isDebug bool, data Data) {
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
func (qb QB) BindModel(model Model) QB {
	valuePtr := reflect.ValueOf(model)
	value := valuePtr.Elem()
	valueType := value.Type()
	if qb.Table == "" {
		tableName := valuePtr.MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
		if qb.Table == "" {
			panic(errors.New("tableName is empty string"))
		}
	}
	if structField, has := valueType.FieldByName("DeletedAt"); has {
		qb.SoftDelete = Column(structField.Tag.Get("db"))
	}
	return qb
}


func (qb QB) Paging(page int , perPage int) QB {
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 10
		log.Print("gofree: Paging(page, perPage) alert perPage is 0 ,perPage can't be 0 . gofree will set perPage 10. but you nedd check your code.")
	}
	qb.Offset = (page - 1) * perPage
	qb.Limit = perPage
	return qb
}
