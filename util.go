package f
import (
	"github.com/google/uuid"
	"reflect"
)
type Data map[Column]interface{}
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
			selectList = append(selectList, dbTag)
		}
		qb.Select = StringsToColumns(selectList)
	}
}
func StringsToColumns(strings []string) (columns []Column) {
	for _, s := range strings {
		columns = append(columns, Column(s))
	}
	return
}
func ColumnsToStrings (columns []Column) (strings []string) {
	for _, column := range columns {
		strings = append(strings, string(column))
	}
	return
}