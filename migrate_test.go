package f_test

import (
	_ "database/sql"
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	"testing"
)
func Migrate202001011113CreateUserTable(mi f.Migrate) {
	mi.CreateTable(f.CreateTableInfo{
		TableName: "user",
		Fields: append(
			[]f.MigrateField{
				mi.Field("id").Char(36).PrimaryKey(),
				mi.Field("name").Varchar(20),
				mi.Field("gender").Char(20).Commit("male | female | unknown"),
				mi.Field("password").Char(64).Commit("sha256"),
				mi.Field("password_salt").Char(36).Commit("uuid salt"),
			},
			mi.CUDTimestamp()...,
			),
		Engine: mi.Engine().InnoDB,
		Collate: mi.Utf8mb4_unicode_ci(),
	})
}
func TestCreateTable(t *testing.T) {
	as := gtest.NewAS(t)
	_=as
	masterDB, err := NewDB()
	if err != nil {ge.Check(err)}
	masterMi := f.NewMigrate(masterDB)
	Migrate202001011113CreateUserTable(masterMi)
}