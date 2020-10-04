package projectNameMigrate

import f "github.com/og/gofree"

func (MasterMigrate) Migrate20201004213024InsertUserTableInintData (mi f.Migrate){
	mi.Exec("INSERT INTO user (id, name) VALUES(?,?)", f.UUID(), "gofree")
}