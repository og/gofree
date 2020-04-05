package exmapleMigrate

import f "github.com/og/gofree"

func Migrate2020_04_04_16_03_07_addbook(mi f.Migrate) {
	// Do not change or remove this code Start
	/****/ mi.MigrateName("2020_04_04_16_03_07_addbook") /****/
	// Do not change or remove this code End

	mi.Table("user").Modify(mi.Field("name").Varchar(20))
	// ALTER  TABLE `user` MODIFY  COLUMN `name`  varchat(20);

	mi.CreateTable(f.CreateTableInfo{
		TableName: "book",
		Fields: append([]f.MigrateField{
			mi.Field("id").Int(10).Unsigned().AutoIncrement(),
			mi.Field("name").Varchar(255).Collate(mi.Utf8mb4_unicode_ci()),
			mi.Field("batch").Int(11),
			mi.Field("bool").Tinyint(1),
			mi.Field("data").Text().Collate(mi.Utf8mb4_unicode_ci()),
		}, mi.CUDTimestamp()...),
		Engine: mi.InnoDB(),
		DefaultCharset: mi.Utf8mb4(),
		Collate: mi.Utf8mb4_unicode_ci(),
	})
	// CREATE TABLE `gofree_migrations` (
	// 	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	// 	`name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	// 	`batch` int(11) NOT NULL DEFAULT '',
	// 	`data` text COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	//  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// 	`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	// 	`deleted_at` timestamp NULL DEFAULT NULL,
	// 	PRIMARY KEY (`id`)
	// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
}