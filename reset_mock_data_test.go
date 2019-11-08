package f_test

import (
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestResetMockData(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User: "root",
		Password: "password",
		Host: "localhost",
		Port: "3306",
		DB: "test_gofree",
	})
	jsonfilepath := filepath.Join(ge.GetString(filepath.Abs("./")), "test", "mock.json")
	_, err := f.ResetMockData(db , map[string]interface{}{}, jsonfilepath) ; ge.Check(err)
	type Mock struct {
		ID int `db:"id"`
		Name string `db:"name"`
	}
	mockList := []Mock{}
	ge.Check(db.Core.Select(&mockList, `SELECT * FROM mock`))
	assert.Equal(t, mockList, []Mock{{1,"nimo"},{2,"nico"}})
}
