package exmapleMigrate

import f "github.com/og/gofree"

func GetProjectDB() f.Database {
	return f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User:       "root",
		Password:   "somepass",
		Host:       "localhost",
		Port:       "3306",
		DB:         "test_gofree",
	})
}
func Migrate() {
	masterDB := GetProjectDB()
	defer masterDB.Close()
	migrate := f.NewMigrate(masterDB)
	migrate.Init(masterDB)
	Migrate2020_04_04_16_03_07_addbook(migrate)
}