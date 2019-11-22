package f_test

import (
	"database/sql"
	_ "database/sql"
	f "github.com/og/gofree"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	ID int `db:"id"`
	Name string `db:"name"`
	IsSuper bool `db:"is_super"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
func (user User) TableName() string {
	return "user"
}
func TestNewDatabase(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User: "root",
		Password: "password",
		Host: "localhost",
		Port: "3306",
		DB: "test_gofree",
	})
	{
		count := db.CountQB(&User{}, f.QB{
			Where: f.And("id", 1),
		})
		assert.Equal(t, count ,1)
	}
	{
		count := db.CountQB(&User{}, f.QB{
			Where: f.And("id", 3),
		})
		assert.Equal(t, count ,0)
	}

	{
		user := User{}
		has := false
		db.OneQB(&user, &has, f.QB{
			Where: f.And("id", "1"),
		})
		assert.Equal(t, user, User{
			ID: 1,
			Name: "nimo",
			IsSuper: false,
			DeletedAt: sql.NullTime{},
		})
		assert.Equal(t, has, true)
	}
	{
		user := User{}
		has := false
		db.OneQB(&user, &has, f.QB{
			Where: f.And("id", "3"),
		})
		assert.Equal(t, user, User{
			ID: 9,
			Name: "",
			IsSuper: false,
			DeletedAt: sql.NullTime{},
		})
		assert.Equal(t, has, false)
	}

	{
		user := User{}
		has := false
		db.OneID(&user, &has, "1")
		assert.Equal(t, user, User{
			ID: 1,
			Name: "nimo",
			IsSuper: false,
			DeletedAt: sql.NullTime{},
		})
		assert.Equal(t, has, true)
	}
	{
		user := User{}
		has := false
		db.OneID(&user, &has, "3")
		assert.Equal(t, user, User{
			ID: 9,
			Name: "",
			IsSuper: false,
			DeletedAt: sql.NullTime{},
		})
		assert.Equal(t, has, false)
	}

}
