package f_test

import (
	f "github.com/og/gofree"
	grand "github.com/og/x/rand"
	"testing"
)
func TestTx(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "password",
		Host:       "localhost",
		Port:       "3306",
		DB:         "test_gofree",
	})
	tx := db.Tx() ; defer func() { f.EndTx(tx, recover()) }()
	{
		user := User{}
		has := false
		db.TxOneID(tx, &user, &has, "1")
		user.Name = grand.StringLetter(4)
		db.TxUpdate(tx, &user)
	}
	{
		user := User{}
		has := false
		db.TxOneID(tx, &user, &has, "10")
		user.Name = grand.StringLetter(1)
		db.TxUpdate(tx, &user)
	}
}
