package main

import (
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	f "github.com/og/gofree"
	migrateAction "github.com/og/gofree/example/cmd/migrate/action"
	connectRDS "github.com/og/gofree/example/connect/rds"
)


func main () {
	rds, err := connectRDS.NewRDS() ; if err != nil {panic(err)}
	f.ExecMigrate(rds.Main ,&migrateAction.MasterMigrate{})
}
