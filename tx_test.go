package f_test

import (
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	grand "github.com/og/x/rand"
	gtest "github.com/og/x/test"
	"testing"
)
func TestTx(t *testing.T) {
	as := gtest.NewAS(t)
	db ,err := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "somepass",
		Host:       "localhost",
		Port:       "3306",
		DB:         "test_gofree",
	}) ; ge.Check(err)
	tx := db.Tx() ; defer func() { tx.End( recover()) }()
	{
		user := User{}
		has := false
		db.TxCreate(tx, &User{ID: "1"})
		db.TxOneID(tx, &user, &has, "1")
		as.Equal(has, true)
		user.Name = grand.StringLetter(4)
		db.TxUpdate(tx, &user)
		{
			readUser := User{}
			has := false
			db.OneID(&readUser, &has, "1")
		}
	}
	{
		user := User{
			ID :"10",
		}
		has := false
		_, err := db.Core.Exec(`delete from ` + user.TableName() + " where id= ?", "10")
		if err != nil {
			panic(err)
		}
		db.Create(&user)
		db.TxOneID(tx, &user, &has, "10")
		user.Name = grand.StringLetter(1)
		db.TxUpdate(tx, &user)
	}
	// if (some) {
	// 	tx.Rollback()
	// }
}
