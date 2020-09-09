package exampleGofree

import (
	f "github.com/og/gofree"
)

var DataSourceName = f.DataSourceName{
	DriverName: "mysql",
	User:       "root",
	Password:   "password",
	Host:       "localhost",
	Port:       "3306",
	DB:         "test_gofree",
}
