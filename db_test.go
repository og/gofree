package f_test

import (
	_ "database/sql"
	f "github.com/og/gofree"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func NewDB() f.Database {
	return f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "password",
		Host:       "localhost",
		Port:       "3306",
		DB:         "test_gofree",
	})
	{
		count := db.CountQB(&User{}, f.QB{
			Where: f.And("id", 1),
		})
		assert.Equal(t, count, 1)
	}
	{
		count := db.CountQB(&User{}, f.QB{
			Where: f.And("id", -1),
		})
		assert.Equal(t, count, 0)
	}

	{
		user := User{}
		has := false
		db.OneQB(&user, &has, f.QB{
			Where: f.And("id", 1),
		})
		assert.Equal(t, user.ID, "1")
		assert.Equal(t, user.Name, "nimo")
		assert.Equal(t, user.IsSuper, false)
		assert.Equal(t, has, true)
	}
	{
		user := User{}
		has := false
		db.OneQB(&user, &has, f.QB{
			Where: f.And("id", -1),
		})
		assert.Equal(t, user.ID, "")
		assert.Equal(t, user.Name, "")
		assert.Equal(t, user.IsSuper, false)
		assert.Equal(t, has, false)
	}

	{
		user := User{}
		has := false
		db.OneID(&user, &has, "1")
		assert.Equal(t, user.ID, "1")
		assert.Equal(t, user.Name, "nimo")
		assert.Equal(t, user.IsSuper, false)
		assert.Equal(t, has, true)
	}
	{
		user := User{}
		has := false
		db.OneID(&user, &has, "-1")
		assert.Equal(t, user.ID, "")
		assert.Equal(t, user.Name, "")
		assert.Equal(t, user.IsSuper, false)
		assert.Equal(t, has, false)
	}
	{
		userList := []User{}
		db.ListQB(&userList, f.QB{
			Where: f.And("id", f.In([]string{"1","2"})),
		})
		assert.Equal(t, userList[0].ID, "1")
		assert.Equal(t, userList[0].Name, "nimo")
		assert.Equal(t, userList[0].IsSuper, false)

		assert.Equal(t, userList[1].ID, "2")
		assert.Equal(t, userList[1].Name, "nico")
		assert.Equal(t, userList[1].IsSuper, true)
	}
	{
		user := User{
			Name: "nimo",
		}
		db.Create(&user)
	}
}
