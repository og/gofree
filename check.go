package f

import (
	gjson "github.com/og/x/json"
	"log"
)

type After struct {
	ActualSQL []string
	ExpectSQL []string
}
func (after After) Check(sqls ...string) {
	if len(sqls) == 0 {
		log.Print("gofree: may be you forget Check() arg")
		return
	}
	after.ExpectSQL = append(after.ExpectSQL, sqls...)
	pass := false
	for _, sql := range sqls {
		for _, actualSQL := range after.ActualSQL {
			if actualSQL == sql {
				pass  = true
			}
		}
	}
	if !pass {
		log.Print("SQL check fail:" + gjson.StringUnfold(after))
	}
}
