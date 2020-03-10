package f

import (
	"fmt"
	gconv "github.com/og/x/conv"
	ge "github.com/og/x/error"
	gjson "github.com/og/x/json"
	glist "github.com/og/x/list"
	l "github.com/og/x/log"
	gmap "github.com/og/x/map"
	grand "github.com/og/x/rand"
	gtime "github.com/og/x/time"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)
type fileExistKind struct {
	Exist bool
	NotExist bool
	IsDir bool
}
func fileExist(filepath string) (kind fileExistKind)  {
	fileInfo, err := os.Stat(filepath)
	if !os.IsExist(err) {
		kind.NotExist = true
		return
	}
	if fileInfo.IsDir() {
		kind.IsDir = true
		return
	}
	kind.Exist = true
	return
}
type mockJSONItem struct {
	Data []map[string]interface{} `json:"data"`
}
type mockDataLocal map[string]interface{}

func (local mockDataLocal) String(key string) string {
	value, has := local[key].(string)
	if !has {panic("Mock.Local  " + key + " not found")}
	return value
}
func (local mockDataLocal) Int(key string) int {
	value, has := local[key].(int)
	if !has {panic("Mock.Local  " + key + " not found")}
	return value
}
func (local mockDataLocal) Float64(key string) float64 {
	value, has := local[key].(float64)
	if !has {panic("Mock.Local  " + key + " not found")}
	return value
}
func (local mockDataLocal) Bool(key string) bool {
	value, has := local[key].(bool)
	if !has {panic("Mock.Local  " + key + " not found")}
	return value
}
type MockData struct {
	Local mockDataLocal `json:"local"`
	Table map[string][]map[string]interface{}
}

func mockFillHelperFunc (mock MockData, helpers map[string]interface{} ,value interface{}) (result string, match bool) {
	valueString := value.(string)
	if strings.HasPrefix(valueString, "@local.") {

		reg, err := regexp.Compile(`\@local\.(.+)\@`) ; ge.Check(err)
		localKey := reg.ReplaceAllString(valueString, "$1")
		localValue, has := mock.Local[localKey]
		if !has {
			panic("localKey " + localKey + "not found")
		}
		return localValue.(string), true
	}
	for varKey, varValue := range helpers {
		varValueReflect := reflect.ValueOf(varValue)
		if valueString == "@" + varKey + "@" {
			if varValueReflect.Kind() == reflect.Func {
				resultList := varValueReflect.Call([]reflect.Value{})
				return resultList[0].String(), true
			}
		}
	}
	if strings.HasPrefix(valueString, "@") && strings.HasSuffix(valueString, "@") {
		reg, err:= regexp.Compile(`^@(.*\()(.+)(\).*)@$`); ge.Check(err)
		funcKey := reg.ReplaceAllString(valueString, "$1$3")
		funcArg := reg.ReplaceAllString(valueString, "$2")
		if funcArg == "" {
			return
		}
		var any []interface{}
		argJSON := "["+ funcArg + "]"
		parserErr := gjson.ParseWithErr(argJSON, &any)
		if parserErr !=nil {
			l.V("parser funcArg error: ", valueString, argJSON)
			panic(err)
		}
		var reflectList []reflect.Value
		for _, item := range any {
			reflectList = append(reflectList, reflect.ValueOf(item))
		}
		defer func() {
			err := recover()
			if err !=nil {
				panic(errors.New("funcKey " + funcKey + " not found"))
			}
		}()
		resultList := reflect.ValueOf(helpers[funcKey]).Call(reflectList)
		result = resultList[0].String()
		return result, true
	}
	return
}

var mockDatahelper = Map{
	"uuid()": func() string {
		return UUID()
	},
	"fromToday()": func(arg ...interface{}) string {
		dayDiff := ge.GetInt(gconv.StringInt(fmt.Sprintf("%v", arg[0])))
		return time.Now().AddDate(0, 0, dayDiff).Format(gtime.Day + " 00:00:00")
	},
	"formTodayHMS()": func(arg ...interface{}) string {
		dayDiff := ge.GetInt(gconv.StringInt(fmt.Sprintf("%v", arg[0])))
		hms := fmt.Sprintf("%s", arg[1])
		return time.Now().AddDate(0, 0, dayDiff).Format(gtime.Day) + " " + hms
	},
	"letter()": func(v interface{}) string {
		size := ge.GetInt(gconv.StringInt(fmt.Sprintf("%v", v)))
		return grand.StringLetter(size)
	},
	"openid()": func() string {
		return grand.StringLetter(28)
	},
	"unionid()": func() string {
		return grand.StringLetter(29)
	},
}

func ResetMockData(db Database, customHelpers map[string]interface{}, filepath string) (mock MockData, err error) {
	helpers := map[string]interface{}{}
	for key, value := range mockDatahelper {
		helpers[key] = value
	}
	for key, value := range customHelpers {
		if !strings.HasPrefix(key, "x-") {
			panic(errors.New("customHelpers: key must have x-"))
		}
		helpers[key] = value
	}

	if !strings.HasPrefix(db.GetDataSourceName().DB, "test_") {
		return mock, errors.New("Danger!!! you try CreateMockData in " + db.GetDataSourceName().DB + ", your code may have errors! (CreateMockData only insert  name like 'test_databasename' database )")
	}
	byteList, err := ioutil.ReadFile(filepath)
	if err != nil { return }
	gjson.ParseByte(byteList, &mock)
	if mock.Local == nil {
		mock.Local = map[string]interface{}{}
	}
	for tableName, item := range mock.Table {
		_, err = db.Core.Exec("truncate table `" + tableName + "`"); ge.Check(err)
		for _, row := range item {
			isCommentRow := glist.StringList{gmap.Keys(row).String()}.Every(func(index int, key string) bool {
				if strings.HasPrefix(key, "#") {
					delete(row, key)
				}
				return key == "//"
			})
			if isCommentRow {
				continue
			}
			for localKey, localValue := range mock.Local {
				if reflect.TypeOf(localValue).String() == "string" {
					valueString := localValue.(string)
					if strings.HasPrefix(valueString, "@") && strings.HasSuffix(valueString, "@") {
						result , match := mockFillHelperFunc(mock, helpers, valueString)
						if match {
							mock.Local[localKey] = result
							continue
						}
					}
				}
			}
			for key, value := range row {
				var repalced bool
				if repalced {continue}
				if repalced {
					continue
				}
				if reflect.ValueOf(value).Kind().String() != "string" {
					continue
				}
				result, match := mockFillHelperFunc(mock, helpers, value)
				if match {
					row[key] = result
					continue
				}
			}
			insertSQL, values := QB{
				Table: tableName,
				Insert: row,
			}.GetInsert()
			_, err := db.Core.Exec(insertSQL, values...)
			if err != nil {
				l.V(insertSQL, values)
				panic(err)
			}
		}
	}
	return
}
