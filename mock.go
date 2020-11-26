package f

import (
	"context"
	"errors"
	"fmt"
	"github.com/og/x/error"
	"reflect"
	"regexp"
)


type Mock struct {
	Tables []interface{}
}
const resetDangerWaring = `
Danger       Danger           Danger          Danger
gofree: ReserAndMock(db,mock) db.DB name prefix must be "test_"
Maybe you reset production environment database
db:
	`
func ResetAndMock(db Database, mock Mock) {
	dataSourceName := db.DataSourceName()
	if !ge.Bool(regexp.MatchString("^test_", dataSourceName.DB)) {
		panic(errors.New(resetDangerWaring + dataSourceName.DB))
	}
	existTruncateTableName := map[string]bool{}
	for _, table := range mock.Tables{
		tableValue := checkAndReturnTable(table)
		mockTruncateTable(db, tableValue, existTruncateTableName)
		// truncate table
		for index := 0;index< tableValue.Len();index ++ {
			model := tableValue.Index(index).Addr().Interface()
			ge.Check(db.Create(context.TODO(), model.(Model)))
		}
	}
}
func checkAndReturnTable(table interface{}) (tableValue reflect.Value){
	tableValue = reflect.ValueOf(table)
	if tableValue.Kind() != reflect.Slice {
		panic(errors.New("mock.Tables must be slice"))
	}
	if tableValue.Len() == 0 {
		panic(errors.New("mock.Tables item can not be nil slice:" + fmt.Sprintf("%+v", tableValue.Interface())))
	}
	return
}
func mockTruncateTable(db Database, tableValue reflect.Value, existTruncateTableName map[string]bool) {

		tableName := ""
		output := tableValue.Index(0).MethodByName("TableName").Call([]reflect.Value{})
		tableName = output[0].String()
		if _, exist := existTruncateTableName[tableName]; exist {
			panic(errors.New("mock.Tables tableName:" + tableName + " already reset and mock"))
		}
		existTruncateTableName[tableName] = true
		_, err := db.Core.Exec("truncate table `" + tableName + "`");
		ge.Check(err)
}