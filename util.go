package f
import (
	"github.com/google/uuid"
	ge "github.com/og/x/error"
	"reflect"
	"regexp"
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