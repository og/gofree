package f

import (
glist "github.com/og/x/list"
gmap "github.com/og/x/map"
"strings"
)
// 通过 DataSourceName 结构体创建 sq.Open(driverName, dataSourceName string) 中的 dataSourceName 字符串
// 因为自己写 root:password@(localhost:3306)/test_gofree 这样的字符串太烦了还容易错

// DataSourceName.User 是个 map[string]stirng 结构
// 用于生成 dataSourceName 中的 ?charset=utf8&parseTime=True&loc=Local 部分
type DataSourceName struct {
	DriverName string
	User string
	Password string
	Host string
	Port string
	DB string
	Query map[string]string
}
func (config DataSourceName) GetString() (dataSourceName string) {
	configList := []string{
		config.User,
		":",
		config.Password,
		"@",
		"(",
		config.Host,
		":",
		config.Port,
		")",
		"/",
		config.DB,
		"?",
	}
	if len(config.Query) == 0 {
		config.Query = map[string]string{
			"charset": "utf8",
			"parseTime": "True",
			"loc": "Local",
		}
	}

	configList = append(configList)
	var UserList glist.StringList
	for _, key := range gmap.StringStringKeys(config.Query) {
		value := config.Query[key]
		UserList.Push(key +"="+value)
	}
	dataSourceName = strings.Join(configList,"") + UserList.Join("&")
	return
}

