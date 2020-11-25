package f_test

import (
	"context"
	"errors"
	f "github.com/og/gofree"
	gconv "github.com/og/x/conv"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	"testing"
	"time"
)

var db f.Database
func init () {
	var err error
	db , err = f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "somepass",
		Host:       "127.0.0.1",
		Port:       "3306",
		DB:         "test_gofree",
	})
	if err != nil {panic(err)}
}
type MasterMigrate struct {

}
func (MasterMigrate) Migrate20201013140601CreateUser(mi f.Migrate) {
	mi.CreateTable(f.CreateTableQB{
		TableName: "user",
		PrimaryKey: "id",
		Fields: append([]f.MigrateField{
			mi.Field("id").Char(36).DefaultString(""),
			mi.Field("name").Varchar(20).DefaultString(""),
			mi.Field("age").Int(11).DefaultInt(0),
			mi.Field("is_super").Tinyint(1).DefaultInt(0),
		}, mi.CUDTimestamp()...),
		Key: nil,
		Engine: mi.Engine().InnoDB,
		Charset: mi.Charset().Utf8mb4,
		Collate: mi.Utf8mb4_unicode_ci(),
	})
}
func (MasterMigrate) Migrate20201016140601CreateUser(mi f.Migrate) {
	mi.CreateTable(f.CreateTableQB{
		TableName: "log",
		PrimaryKey: "id",
		Fields: []f.MigrateField{
			mi.Field("id").Int(11).AutoIncrement(),
			mi.Field("message").Varchar(20).DefaultString(""),
		},
		Key: nil,
		Engine: mi.Engine().InnoDB,
		Charset: mi.Charset().Utf8mb4,
		Collate: mi.Utf8mb4_unicode_ci(),
	})
}
func (MasterMigrate) Migrate20201102152223CreateUserLocation(mi f.Migrate) {
	mi.CreateTable(f.CreateTableQB{
		TableName: "user_location",
		PrimaryKey: "id",
		Fields: []f.MigrateField{
			mi.Field("id").Char(36),
			mi.Field("point").Type("POINT", 0),
		},
		Key: nil,
		Engine: mi.Engine().InnoDB,
		Charset: mi.Charset().Utf8mb4,
		Collate: mi.Utf8mb4_unicode_ci(),
	})
}


func TestDB(t *testing.T) {
	as := gtest.NewAS(t)
	f.ExecMigrate(db, &MasterMigrate{})
	_, err := db.DB.Exec(`truncate table user`) ; ge.Check(err)
	db.Create(&User{
		Name:      "nimo",
		Age:       18,
		IsSuper:   true,
	})
	var user User
	var hasUser bool
	db.OneQB(&user, &hasUser, f.QB{
		Where:      f.And(user.Column().Name, "nimo"),
		Check: []string{"SELECT `id`, `name`, `age`, `is_super`, `created_at`, `updated_at`, `deleted_at` FROM `user` WHERE `name` = ? AND `deleted_at` IS NULL LIMIT ?"},
	})
	as.Equal(len(user.ID), 36)
	as.Equal(user.Name, "nimo")
	as.Equal(user.Age, 18)
	as.Equal(user.IsSuper , true)
	as.True(user.CreatedAt.After(time.Now().Add(-3*time.Second)))
	as.True(user.CreatedAt.Before(time.Now().Add(time.Second)))
	as.True(user.UpdatedAt.After(time.Now().Add(-3*time.Second)))
	as.True(user.UpdatedAt.Before(time.Now().Add(time.Second)))
	{
		ctx ,_ := context.WithTimeout(context.Background(), time.Nanosecond)
		err = db.CoreCreate(f.SqlOpt{
			Context: ctx,
		}, &User{Name:"gofree"})
		as.Error(err, context.DeadlineExceeded)
		var createError error
		select {
		case <-ctx.Done():
			createError = ctx.Err()
		default:
			t.Fatal("Should not be pass")
		}
		as.Error(createError, context.DeadlineExceeded)
	}
	{
		err = db.Tx(func(tx *f.Tx) error {
			ge.Check(db.CoreCreate(f.SqlOpt{
				Tx: tx,
			}, &User{
				Name: "TXFAIL",
			}))
			var hasUser bool
			err = db.CoreOneQB(f.SqlOpt{Tx:tx,}, &user, &hasUser, f.QB{
				Where:      f.And(user.Column().Name, "TXFAIL"),
			})
			as.NoError(err)
			as.True(hasUser)
			as.Equal(user.Name , "TXFAIL")
			return errors.New("some error")
		})
		as.ErrorString(err, "some error")
		user := User{}
		var hasUser bool
		db.OneQB(&user,&hasUser, f.QB{Where: f.And(user.Column().Name, "TXFAIL")})
		as.Equal(hasUser, false)
		as.Equal(user.ID, IDUser(""))
	}
	{
		err = db.Tx(func(tx *f.Tx) error {
			ge.Check(db.CoreCreate(f.SqlOpt{
				Tx: tx,
			}, &User{
				Name: "TXCOMMIT",
			}))
			return nil
		})
		as.NoError(err)
		user := User{}
		var hasUser bool
		db.OneQB(&user,&hasUser, f.QB{Where: f.And(user.Column().Name, "TXCOMMIT")})
		as.Equal(hasUser, true)
		as.Equal(user.Name, "TXCOMMIT")
	}
	{
		newLog := Log{Message: "abc"}
		err := db.CoreCreate(f.SqlOpt{}, &newLog)
		as.NoError(err)
		as.Gt(int(newLog.ID), 0)
	}
	{
		newLog := Log2{Message: "abc"}
		err := db.CoreCreate(f.SqlOpt{}, &newLog)
		as.NoError(err)
		as.Gt(int(newLog.ID), 0)
	}
	{
		newLog := Log3{Message: "abc"}
		err := db.CoreCreate(f.SqlOpt{}, &newLog)
		as.NoError(err)
		idInt, err := gconv.StringInt(newLog.ID)
		as.NoError(err)
		as.Gt(idInt, 0)
	}
	{
		newLog := Log4{Message: "abc"}
		as.PanicError("Log4.ID type must be uint or int or string", func() {
			err := db.CoreCreate(f.SqlOpt{}, &newLog)
			if err != nil {panic(err)}
		})
	}
	{
		newLog := Log5{Message: "abc"}
		as.PanicError(`dbAutoIncrement muse be dbAutoIncrement:"true" or dbAutoIncrement:"false" can not be dbAutoIncrement:"t"`, func() {
			err := db.CoreCreate(f.SqlOpt{}, &newLog)
			if err != nil {panic(err)}
		})
	}

}