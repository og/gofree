package f_test

import (
	f "github.com/og/gofree"
	gtest "github.com/og/x/test"
	"testing"
)

func TestCreateTableQB_ToSQL(t *testing.T) {
	as := gtest.NewAS(t)
	mi := f.Migrate{}
	qb := f.CreateTableQB{
		TableName: "user",
		PrimaryKey: []string{"id","order_id"},
		Fields: append([]f.MigrateField{
			mi.Field("id").Char(36).DefaultString(""),
			mi.Field("name").Varchar(20).DefaultString(""),
			mi.Field("age").Int(11),
			mi.Field("disabled").Tinyint(1),
			mi.FieldRaw("`is_super` tinyint(4) NOT NULL"),
		}, mi.CUDTimestamp()...),
		Engine: mi.Engine().InnoDB,
		Charset: mi.Charset().Utf8mb4,
		Collate: mi.Utf8mb4_unicode_ci(),
	}
	_=qb
	_=as
}