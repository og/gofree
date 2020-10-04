package exampleGofree

import (
	f "github.com/og/gofree"
)


func NewDB() (db f.Database, err error) {
	dataSourceName := f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "somepass",
		Host:       "localhost",
		Port:       "3306",
		DB:         "example_gofree",
	}
	return f.NewDatabase(dataSourceName)
}