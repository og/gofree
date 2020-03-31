package f_test

import (
	"database/sql"
	_ "database/sql"
	"errors"
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	gtime "github.com/og/x/time"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
func init () {
	_= errors.New
}

type User struct {
	ID string `db:"id"`
	Name string `db:"name"`
	IsSuper bool `db:"is_super"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
const fUserID = "id"
const fUserName = "name"
const fUserIsSuper = "is_super"
const fUserCreatedAt = "created_at"
const fUserUpdatedAt = "updated_at"
const fUserDeletedAt = "deleted_at"
func (User) TableName() string {
	return "user"
}
func (model *User) BeforeCreate() {
	if model.ID == "" {
		model.ID = f.UUID()
	}
}
func TestNewDatabase(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "somepass",
		Host:       "localhost",
		Port:       "3306",
		DB:         "test_gofree",
	})
	{
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal(errors.New("sholud be error"))
				}

				err := r.(error)

				assert.Equal(t, "db.Update(&model) or db.TxUpdate(&model) model.id is zero", err.Error())
			}()
			user := User{}
			db.Update(&user)
		}()
	}
	{
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal(errors.New("sholud be error"))
				}

				err := r.(error)

				assert.Equal(t, "db.Update(&model) or db.TxUpdate(&model) model.id is zero", err.Error())
			}()
			user := User{}
			db.Update(&user)
		}()
	}
	{
		user := User{
			ID: "1",
		}
		_, err := db.Core.Exec(`delete from ` + user.TableName() + " where id= ?", "1") ; ge.Check(err)
		user.Name = "nimo"
		user.IsSuper = false
		db.Create(&user)
	}
	{
		count := db.CountQB(&User{}, f.QB{
			Where: f.And("id", "1"),
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
		_, err := db.Core.Exec(`DELETE FROM user WHERE id = ?`, "2")
		if err !=nil {
			panic(err)
		}
		db.Create(&User{
			ID: "2",
			Name: "nico",
			IsSuper: true,
		})
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
	{
		user := User{
			Name: "deletedQB",
		}
		db.Create(&user)
		db.DeleteQB(&user, f.QB{
			Where: f.And("id", user.ID),
		})
		userList := []User{}
		query, values := f.QB{
			Table: user.TableName(),
			Where: f.And("id", user.ID),
		}.GetSelect()
		err := db.Core.Select(&userList, query, values...) ; ge.Check(err)
		assert.Equal(t, len(userList), 1)
		assert.Equal(t, userList[0].ID, user.ID)
		assert.Equal(t, userList[0].Name, "deletedQB")
		assert.Equal(t, userList[0].DeletedAt.Valid, true)
	}
	{
		user := User{
			Name: "deleted",
		}
		db.Create(&user)
		db.Delete(&user)
		userList := []User{}
		query, values := f.QB{
			Table: user.TableName(),
			Where: f.And("id", user.ID),
		}.GetSelect()
		err := db.Core.Select(&userList, query, values...) ; ge.Check(err)
		assert.Equal(t, len(userList), 1)
		assert.Equal(t, userList[0].ID, user.ID)
		assert.Equal(t, userList[0].Name, "deleted")
		assert.Equal(t, userList[0].DeletedAt.Valid, true)
	}
	{
		user := User{
			Name: "update1",
		}
		db.Create(&user)
		assert.Equal(t, len(user.ID), 36)
		assert.Equal(t, user.Name, "update1")
		user.Name = "update2"
		db.Update(&user)
		assert.Equal(t, user.Name, "update2")
		userList := []User{}
		query, values := f.QB{
			Table: user.TableName(),
			Where: f.And("id", user.ID),
		}.GetSelect()
		err := db.Core.Select(&userList, query, values...) ; ge.Check(err)
		assert.Equal(t, len(userList), 1)
		assert.Equal(t, userList[0].ID, user.ID)
		assert.Equal(t, userList[0].Name, "update2")
		assert.Equal(t, userList[0].UpdatedAt.Format(gtime.Minute), time.Now().Format(gtime.Minute))
	}
}

type User2 struct {
	EventID string
	ID string `db:"id"`
	Name string `db:"name"`
	IsSuper bool `db:"is_super"`
}
func (User2) TableName() string {
	return "user"
}
func (model *User2) BeforeCreate(){
	if model.ID == "" {
		model.ID = f.UUID()
	}
}
func TestCreateIgnoreField(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "somepass",
		Host:       "localhost",
		Port:       "3306",
		DB:         "test_gofree",
	})
	user2 := User2{
		ID: "user2",
	}
	_, err := db.Core.Exec(`delete from ` + user2.TableName() + " where id= ?", "user2") ; ge.Check(err)
	db.Create(&user2)
	newUser := User2{}
	has := false
	db.OneID(&newUser, &has, "user2")
	assert.Equal(t, has, true)
	assert.Equal(t, newUser.ID, "user2")
	{
		// 正常情况下这个不应该报错
		db.ListQB(&[]User2{}, f.QB{
			Check: []string{"SELECT `id`, `name`, `is_super` FROM `user` WHERE `deleted_at` IS NULL"},
		})
	}
}