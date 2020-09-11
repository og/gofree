package f

import (
)
type Migrate struct {
	db Database
}

type MigrateModel struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Batch int `db:"batch"`
	Data string `db:"data"`
}
const createMigrateSQL = `
CREATE TABLE  IF NOT EXISTS gofree_migrations (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  batch int(11) NOT NULL,
  data text COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 
`
func (mi Migrate) Init(db Database) {
	_, err := db.Core.Exec(createMigrateSQL)
	if err != nil {
		panic(err)
	}
}
func NewMigrate (db Database) Migrate {
	return Migrate{
		db: db,
	}
}
type MigrateEngine string
type MigrateCharset string
type MigrateCollate string
type CreateTableInfo struct {
	TableName string
	Fields []MigrateField
	Engine MigrateEngine
	DefaultCharset MigrateCharset
	Collate MigrateCollate
}
type MigrateField struct {
	name string
	size int
	fieldType string
	unsigned bool
	null bool
	autoIncrement bool
	callate string
	defaultValue migrateDefaultValue
	primaryKey bool
	references struct{
		valid bool
		otherTableName string
		otherTableField string
	}
	unique bool
	extra []string
	commit string
	raw string
}
func (mi MigrateField) Int(size int) MigrateField {
	mi.size = size
	mi.fieldType = "int"
	return mi
}
func (mi MigrateField) Char(size int) MigrateField {
	mi.size = size
	mi.fieldType = "char"
	return mi
}
func (mi MigrateField) Varchar(size int) MigrateField {
	mi.size = size
	mi.fieldType = "varchar"
	return mi
}
func (mi MigrateField) Unsigned() MigrateField {
	mi.unsigned = true
	return mi
}
func (mi Migrate) Utf8mb4_unicode_ci () MigrateCollate {
	return "utf8mb4_unicode_ci"
}
func (mi Migrate) Engine() (e struct {
	BLACKHOLE MigrateEngine
	CSV MigrateEngine
	InnoDB MigrateEngine
	MEMORY MigrateEngine
	MRG_MyISAM MigrateEngine
	MyISAM MigrateEngine
	PERFORMANCE_SCHEMA MigrateEngine
}) {
	e.BLACKHOLE = "BLACKHOLE"
	e.CSV = "CSV"
	e.InnoDB = "InnoDB"
	e.MEMORY = "MEMORY"
	e.MRG_MyISAM = "MRG_MyISAM"
	e.MyISAM = "MyISAM"
	e.PERFORMANCE_SCHEMA = "PERFORMANCE_SCHEMA"
	return
}
func (mi Migrate) Charset() (v struct {
	Utf8mb4 MigrateCharset

}) {
	v.Utf8mb4 = "utf8mb4"
	return
}

// utf8mb4_unicode_ci

func (mi MigrateField) Collate(kind string)  MigrateField{
	mi.callate = kind
	return mi
}
type migrateDefaultValue struct {
	raw string
}
func (mi Migrate) CurrentTimeStamp() migrateDefaultValue {
	return migrateDefaultValue{
		raw: "CURRENT_TIMESTAMP",
	}
}
func (mi Migrate) DefaultString(s string) migrateDefaultValue {
	return migrateDefaultValue{
		raw: `"` + s + `"`,
	}
}
func (mi MigrateField) Default(value migrateDefaultValue) MigrateField {
	mi.defaultValue = value
	return mi
}
func (mi MigrateField) Null()  MigrateField{
	mi.null = true
	return mi
}
func (mi MigrateField) AutoIncrement() MigrateField {
	mi.autoIncrement = true
	return mi
}

func (mi MigrateField) Text() MigrateField {
	mi.fieldType = "text"
	return mi
}
func (Migrate) MigrateName(name string){}
func (Migrate) CreateTable(info CreateTableInfo) {
	sql := stringQueue{}
	sql.Push("CREATE TABLE `", info.TableName , "`(")
}
type Alter struct {
	migrateField MigrateField
	tableName string
}
func (al Alter) Modify(migrateField MigrateField) Alter {
	al.migrateField = migrateField
	return al
}
func (Migrate) AlterTable(tableName string) Alter {
	return Alter {
		tableName: tableName,
	}
}
func (Migrate) Field(name string) MigrateField {
	return MigrateField{
		name: name,
	}
}
func (mi MigrateField) PrimaryKey() MigrateField {
	mi.primaryKey = true
	return mi
}
func (mi MigrateField) References(otherTableName string, otherTableField string) MigrateField {
	mi.references.valid = true
	mi.references.otherTableName = otherTableName
	mi.references.otherTableField = otherTableField
	return mi
}
func (mi MigrateField) Unique() MigrateField {
	mi.unique = true
	return mi
}
func (mi MigrateField) Timestamp() MigrateField {
	mi.fieldType = "timestamp"
	return mi
}
func (mi MigrateField) Commit(commit string) MigrateField {
	mi.commit = commit
	return mi
}
func (mi MigrateField) Raw(raw string) MigrateField {
	mi.raw = raw
	return mi
}
func (mi MigrateField) Extra(extra string) MigrateField {
	mi.extra = append(mi.extra, extra)
	return  mi
}
func (mi MigrateField) OnUpdateCurrentTimeStamp() MigrateField {
	mi.extra = append(mi.extra, "ON UPDATE CURRENT_TIMESTAMP")
	return mi
}
func (mi Migrate) CreatedAtTimestamp() MigrateField {
	return mi.Field("created_at").
		Timestamp().
		Default(mi.CurrentTimeStamp())
}
func (mi Migrate) UpdatedAtTimestamp() MigrateField {
	return mi.Field("updated_at").
		Timestamp().
		Default(mi.CurrentTimeStamp()).
		OnUpdateCurrentTimeStamp()
}
func (mi Migrate) DeletedAtTimestamp() MigrateField {
	return mi.Field("deleted_at").
		Timestamp().
		Null()
}
func (mi Migrate) CUDTimestamp() []MigrateField {
	return []MigrateField{
		mi.CreatedAtTimestamp(),
		mi.UpdatedAtTimestamp(),
		mi.DeletedAtTimestamp(),
	}
}

func (mi MigrateField) Tinyint(size int) MigrateField {
	mi.fieldType = "tinyint"
	return mi
}