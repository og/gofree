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
	ge.Check(err)
}
type MasterMigrate struct {

}
func (MasterMigrate) Migrate20201013140601CreateUser(mi f.Migrate) {
	mi.CreateTable(f.CreateTableQB{
		TableName: "user",
		PrimaryKey: []string{"id"},
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
		PrimaryKey: []string{"id"},
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
		PrimaryKey: []string{"id"},
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


func init () {
	f.ExecMigrate(db, &MasterMigrate{})
}
func TestDB(t *testing.T) {
	as := gtest.NewAS(t)
	_, err := db.Core.Exec(`truncate table user`) ; ge.Check(err)
	ctx := context.Background()
	ge.Check(db.Create(ctx,&User{
		Name:      "nimo",
		Age:       18,
		IsSuper:   true,
	}))
	var user User
	var hasUser bool
	ge.Check(db.OneQB(ctx, &user, &hasUser, f.QB{
		Where:      f.And(user.Column().Name, "nimo"),
		Check: []string{"SELECT `id`, `name`, `age`, `is_super`, `created_at`, `updated_at`, `deleted_at` FROM `user` WHERE `name` = ? AND `deleted_at` IS NULL LIMIT ?"},
	}))
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
		err = db.Create(ctx, &User{Name:"gofree"})
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
		err = db.Transaction(ctx, func(tx *f.Tx) error {
			ge.Check(tx.Create(ctx, &User{
				Name: "TXFAIL",
			}))
			var hasUser bool
			err = tx.OneQB(ctx, &user, &hasUser, f.QB{
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
		ge.Check(db.OneQB(ctx, &user,&hasUser, f.QB{Where: f.And(user.Column().Name, "TXFAIL")}))
		as.Equal(hasUser, false)
		as.Equal(user.ID, IDUser(""))
	}
	{
		err := db.Transaction(ctx, func(tx *f.Tx) error {
			return tx.Create(ctx, &User{
				Name: "TXCOMMIT SUCCESS",
			})
		})
		as.NoError(err)
		user := User{}
		var hasUser bool
		ge.Check(db.OneQB(ctx, &user,&hasUser, f.QB{Where: f.And(user.Column().Name, "TXCOMMIT SUCCESS")}))
		as.Equal(hasUser, true)
		as.Equal(user.Name, "TXCOMMIT SUCCESS")
	}
	{
		panicValue := as.Panic(func() {
			err := db.Transaction(ctx, func(tx *f.Tx) error {
				ge.Check(tx.Create(ctx, &User{
					Name: "TXCOMMIT PANIC",
				}))
				panic(errors.New("test panic"))
			})
			_=err
		})
		err := panicValue.(error)
		as.ErrorString(err, "test panic")
		user := User{}
		var hasUser bool
		ge.Check(db.OneQB(ctx, &user,&hasUser, f.QB{Where: f.And(user.Column().Name, "TXCOMMIT PANIC")}))
		as.Equal(hasUser, false)
	}
	{
		err := db.Transaction(ctx, func(tx *f.Tx) error {
			ge.Check(tx.Create(ctx, &User{
				Name: "TXCOMMIT ROLLBACK",
			}))
			return tx.Rollback()
		})
		as.NoError(err)
		user := User{}
		var hasUser bool
		ge.Check(db.OneQB(ctx, &user,&hasUser, f.QB{Where: f.And(user.Column().Name, "TXCOMMIT ROLLBACK")}))
		as.Equal(hasUser, false)
	}
	{
		newLog := Log{Message: "abc"}
		err := db.Create(ctx, &newLog)
		as.NoError(err)
		as.Gt(int(newLog.ID), 0)
	}
	{
		newLog := Log2{Message: "abc"}
		err := db.Create(ctx, &newLog)
		as.NoError(err)
		as.Gt(int(newLog.ID), 0)
	}
	{
		newLog := Log3{Message: "abc"}
		err := db.Create(ctx, &newLog)
		as.NoError(err)
		idInt, err := gconv.StringInt(newLog.ID)
		as.NoError(err)
		as.Gt(idInt, 0)
	}
	{
		newLog := Log4{Message: "abc"}
		as.PanicError("log ID type must be uint or int or string", func() {
			err := db.Create(ctx, &newLog)
			if err != nil {panic(err)}
		})
	}
	{
		newLog := Log5{Message: "abc"}
		as.PanicError(`dbAutoIncrement muse be dbAutoIncrement:"true" or dbAutoIncrement:"false" can not be dbAutoIncrement:"t"`, func() {
			err := db.Create(ctx, &newLog)
			if err != nil {panic(err)}
		})
	}
}
func TestUpdate(t *testing.T) {
	as := gtest.NewAS(t)
	ctx := context.Background()
	// update
	_, err := db.Core.Exec(`truncate table user`) ; ge.Check(err)
	{
		user := User{
			Name: "updateName",
		}
		userCol :=User{}.Column()
		ge.Check(db.Create(ctx, &user))
		err := db.UpdateData(ctx, &user, f.Data{
			userCol.Name: "updateName2",
			userCol.IsSuper: true,
		}) ; ge.Check(err)
		as.Equal(user.Name, "updateName2")
		as.Equal(user.IsSuper, true)
		{
			newUser := User{}
			var has bool
			ge.Check(db.OneID(ctx, &newUser, &has, user.ID))
			as.Equal(newUser.Name, "updateName2")
			as.Equal(newUser.IsSuper, true)
		}
	}
}
func TestDatabase_ScanRow(t *testing.T) {
	_, err := db.Core.Exec(`truncate table user`) ; ge.Check(err)
	as := gtest.NewAS(t)
	ctx := context.Background()
	_=as
	user := User{
		Name: "scanRow",
	}
	ge.Check(db.Create(ctx, &user))
	userCol := User{}.Column()
	{
		var has bool
		var name string
		ge.Check(db.QueryRowScan(ctx, &has, f.QB{
			Table: user.TableName(),
			Select: []f.Column{userCol.Name},
			Where: f.And(userCol.ID, user.ID),
		}, &name))
		as.Equal(name, "scanRow")
		as.Equal(has, true)
	}
	{
		var has bool
		var name string
		ge.Check(db.QueryRowScan(ctx, &has, f.QB{
			Table: user.TableName(),
			Select: []f.Column{userCol.Name},
			Where: f.And(userCol.ID, ""),
		}, &name))
		as.Equal(name, "")
		as.Equal(has, false)
	}
}