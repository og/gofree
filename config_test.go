package f_test

import (
	"database/sql"
	f "github.com/og/gofree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleCreateDataSourceName() {
	dataSourceName := f.DataSourceName{
		User: "root",
		Password: "password",
		Host: "localhost",
		Port: "3306",
		DB: "test_gofree",
	}.GetString()
	sql.Open("mysql", dataSourceName)

}

func TestCreateDataSourceName(t *testing.T) {
	// ignore User
	assert.Equal(t,
		"root:password@(localhost:3306)/test_gofree?charset=utf8&loc=Local&parseTime=True",
		f.DataSourceName{
			User: "root",
			Password: "password",
			Host: "localhost",
			Port: "3306",
			DB: "test_gofree",
		}.GetString(),
	)
	// use User
	assert.Equal(t,
		"root:password@(localhost:3306)/test_gofree?charset=gb&loc=local",
		f.DataSourceName{
			User: "root",
			Password: "password",
			Host: "localhost",
			Port: "3306",
			DB: "test_gofree",
			Query: map[string]string{
				"charset": "gb",
				"loc": "local",
			},
		}.GetString(),
	)
}
