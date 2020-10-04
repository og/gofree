package f

import (
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
		log.Printf("SQL check fail: %+v", after)
	}
}
