package main

import (
	f "github.com/og/gofree"
	exampleGofree "github.com/og/gofree/example"
	projectNameMigrate "github.com/og/gofree/example/migrate"
)


func main () {
	db, err := exampleGofree.NewDB() ; if err != nil {panic(err)}
	f.ExecMigrate(db ,&projectNameMigrate.MasterMigrate{})
}
