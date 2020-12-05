package connectRDS

import f "github.com/og/gofree"

// 仅作为演示所以将数据库账号存放在版本控制中，
// 日常项目开发环境和生产环境中数据库账号等不加入代码版本控制（git）
var dataSourceName = f.DataSourceName{
	DriverName: "mysql",
	User:       "root",
	Password:   "somepass",
	Host:       "localhost",
	Port:       "3306",
	DB:         "example_gofree",
}