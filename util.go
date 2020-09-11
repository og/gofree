package f
import (
	"github.com/google/uuid"
	ge "github.com/og/x/error"
	"reflect"
	"regexp"
	"strings"
)
type Map map[string]interface{}
func UUID() string {
	return uuid.New().String()
}
func GetUUID32 (uuid string) string {
	separator, err := regexp.Compile(`-`) ; ge.Check(err)
	return separator.ReplaceAllString(uuid, "")
}
func GetUUID36(uuid32 string) string {
	uuid := ""
	for i:=0;i<len(uuid32);i++ {
		uuid += string(uuid32[i])
		switch i {
		case 7, 11, 15, 19:
			uuid += "-"
		}
	}
	return uuid
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
		qb.Select = selectList
	}
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