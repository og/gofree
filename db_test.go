package f_test

import (
	_ "database/sql"
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User: "root",
		Password: "password",
		Host: "localhost",
		Port: "3306",
		DB: "test_gofree",
	})
	query, values := f.QB{
		Table: "user",
		Where: f.And("id", 1),
		Count: true,
	}.GetSelect()
	row, err := db.Core.Queryx(query, values...) ; ge.Check(err)
	defer func() {
		ge.Check(row.Close())
	}()
	var count int
	row.Next()
	ge.Check(row.Scan(&count))
	assert.Equal(t, count ,1)
}
