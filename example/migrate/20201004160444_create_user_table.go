package projectNameMigrate

import (
	f "github.com/og/gofree"
)

func (MasterMigrate) Migrate20201004160444CreateUserTable(mi f.Migrate) {
	mi.CreateTable(f.CreateTableQB{
		TableName: "user",
		PrimaryKey: "id",
		Fields: append([]f.MigrateField{
			mi.Field("id").Char(36).DefaultString(""),
			mi.Field("name").Varchar(20).DefaultString(""),
			mi.Field("age").Int(11).DefaultInt(0),
			mi.Field("disabled").Tinyint(1).DefaultInt(0),
		}, mi.CUDTimestamp()...),
		Key: map[string][]string{
			"name": []string{"name"},
		},
		Engine: mi.Engine().InnoDB,
		Charset: mi.Charset().Utf8mb4,
		Collate: mi.Utf8mb4_unicode_ci(),
	})
}
