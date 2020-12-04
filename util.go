package f
import (
	"github.com/google/uuid"
	"reflect"
	"strings"
)
type Data map[Column]interface{}
type sortData struct {
	Column Column
	Value interface{}
}
func UUID() string {
	return uuid.New().String()
}
// 解决数据库字段新增,但是没有通知go导致的解析问题
func scanModelMakeSQLSelect(modelType reflect.Type, qb *QB)  {
	if len(qb.Select) == 0 {
		selectList := []string{}
		for i:=0;i<modelType.NumField();i++ {
			dbTag, hasDBTag := modelType.Field(i).Tag.Lookup("db")
			if !hasDBTag {
				continue
			}
			if dbTag == "" {
				continue
			}
			selectList = append(selectList, dbTag)
		}
		qb.Select = StringsToColumns(selectList)
	}
}
func StringsToColumns(strings []string) (columns []Column) {
	for _, s := range strings {
		columns = append(columns, NewColumn(s))
	}
	return
}
func ColumnsToStrings (columns []Column) (strings []string) {
	for _, column := range columns {
		strings = append(strings, column.String())
	}
	return
}

type stringQueue struct {
	Value []string
}
func (v *stringQueue) Push(args... string) {
	v.Value = append(v.Value, args...)
}

func (sList *stringQueue) Pop() stringQueue {
	return sList.PopBind(&stringQueueBindValue{})
}
type stringQueueBindValue struct {
	Value string
	Has bool
}
func (sList *stringQueue) PopBind(last *stringQueueBindValue) stringQueue {
	listLen := len(sList.Value)
	if listLen == 0 {
		/*
			Clear StringListBindValue Because in this case
				```
				list.PopBind(&last)
				// do Something..
				list.PopBind(&last)
				```
				last test same var
		*/
		last.Value = ""
		last.Has = false
		return *sList
	}
	last.Value = sList.Value[listLen-1]
	last.Has = true
	sList.Value = sList.Value[:listLen-1]
	return *sList
}
func (v stringQueue) Join(separator string) string {
	return strings.Join(v.Value, separator)
}

func getPtrElem(ptr interface{}) (value reflect.Value, isPtr bool) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		isPtr = false
		return
	}
	value = v.Elem()
	isPtr = true
	return
}