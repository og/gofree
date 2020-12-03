package f

import (
	"errors"
	ge "github.com/og/x/error"
	"log"
	"reflect"
	"strconv"
	"strings"
)
type Migrate struct {
	db Database
}

const createMigrateSQL = `
CREATE TABLE  IF NOT EXISTS gofree_migrations (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 
`

func (mi Migrate) Init() {
	_, err := mi.db.Core.Exec(createMigrateSQL)
	mi.CheckError(err, createMigrateSQL)
}
func ExecMigrate(db Database, ptr interface{}) {
	rPtrValue := reflect.ValueOf(ptr)
	if rPtrValue.Kind() != reflect.Ptr {
		panic(errors.New("ExecMigrate(db, ptr) ptr must be pointer"))
	}
	rValue := rPtrValue.Elem()
	rType := rValue.Type()
	if rType.PkgPath() == "main" {
		panic(errors.New("ExecMigrate(db, ptr) ptr can not belong to package main"))
	}
	mi := NewMigrate(db)
	mi.Init()
	miValue := reflect.ValueOf(mi)
	methodNames := []string{}
	for i:=0;i<rType.NumMethod();i++ {
		method := rType.Method(i)
		if strings.HasPrefix(method.Name, "Migrate") {
			methodNames = append(methodNames, method.Name)
		}
	}
	for _, methodName := range methodNames {
		row, err := db.Core.Queryx(`SELECT count(*) FROM gofree_migrations WHERE name = ?`, methodName)
		defer row.Close()
		if err != nil {
			panic(err)
		}
		count := 0
		row.Next()
		ge.Check(row.Scan(&count))
		if count == 0 {} else if count == 1 {
			continue
		} else {
			panic(errors.New("warning: gofree_migrations has two same name: " + methodName))
		}
		log.Print("[gofree migrate]exec: " +methodName)
		rValue.MethodByName(methodName).Call([]reflect.Value{miValue})
		_, err = db.Core.Exec("INSERT INTO gofree_migrations (name) VALUES(?)", methodName) ; ge.Check(err)
		log.Printf("[gofree migrate]done: " +methodName)
	}
}
func NewMigrate (db Database) Migrate {
	return Migrate{
		db: db,
	}
}
type MigrateEngine string
func (engine MigrateEngine) String() string {return string(engine)}
type MigrateCharset string
func (charset MigrateCharset) String() string {return string(charset)}
type MigrateCollate string
func (collate MigrateCollate) String() string {return string(collate)}
type CreateTableQB struct {
	TableName string
	PrimaryKey []string
	Fields []MigrateField
	UniqueKey map[string][]string
	Key map[string][]string
	BeforeOfEndBracketRaw []string
	Engine MigrateEngine
	Charset MigrateCharset
	Collate MigrateCollate
}
func (qb CreateTableQB) ToSQL() string {
	sq := stringQueue{}
	if qb.TableName == "" {
		panic(errors.New("TableName can not be empty string"))
	}
	newLine := "\n"
	sq.Push(`CREATE TABLE`, " ", "`", qb.TableName, "`", "(")
	if len(qb.Fields) == 0 {
		panic(errors.New("Fields can not be empty slice"))
	}
	for _, field := range  qb.Fields {
		sq.Push(newLine, "  ")
		if field.raw != "" {
			sq.Push(field.raw, ",")
			continue
		}
		fieldSize := strconv.FormatInt(int64(field.size), 10)
		sq.Push("`", field.name ,"`"," ", field.fieldType)
		if field.size != 0 {
			sq.Push(" (", fieldSize, ")")
		}
		if field.unsigned {
			sq.Push(" unsigned")
		}
		if field.characterSet != "" {
			sq.Push(" CHARACTER SET ", field.characterSet)
		}
		if field.collate != "" {
			sq.Push(" COLLATE ", field.collate)
		}
		if field.null {
			sq.Push(" NULL")
		} else {
			sq.Push(" NOT NULL")
		}
		if field.defaultValue.raw != "" {
			sq.Push(" DEFAULT", " ", field.defaultValue.raw)
		}
		if len(field.extra) != 0 {
			sq.Push(" ")
			sq.Push(strings.Join(field.extra, " "))
		}
		if field.autoIncrement {
			sq.Push(" AUTO_INCREMENT")
		}
		if field.references.valid {
			if field.references.otherTableName == "" {
				panic(errors.New("references tableName can not be empty string"))
			}
			if field.references.otherTableField == "" {
				panic(errors.New("references field can not be empty string"))
			}
			sq.Push(" REFERENCES", field.references.otherTableName, "(", field.references.otherTableField, ")")
		}
		if field.commit != "" {
			sq.Push(" COMMENT", "'" + field.commit + "'")
		}
		sq.Push(",")
	}
	if len(qb.PrimaryKey) == 0 {
		panic(errors.New("your must set PRIMARY KEY "))
	}
	sq.Push(newLine, "  PRIMARY KEY (`", strings.Join(qb.PrimaryKey, "`,`"), "`),")
	for key, values := range qb.UniqueKey {
		sq.Push(newLine, "  UNIQUE KEY ", "`", key, "`", " (`" , strings.Join(values, "`,`") ,"`),")
	}
	for key, values := range qb.Key {
		sq.Push(newLine, "  KEY ", "`", key, "`", " (`" , strings.Join(values, "`,`") ,"`),")
	}
	for _, raw := range qb.BeforeOfEndBracketRaw {
		sq.Push(newLine, strings.TrimSuffix(raw, ","), ",")
	}

	/* 处理sql CRAETE TABLE tableName() 的 () 中不能以 , 结尾的语法  */{
		popValue := stringQueueBindValue{}
		sq.PopBind(&popValue)
		if !popValue.Has {
			panic(errors.New("sq.PopBind() must has value"))
		}
		sq.Push(strings.TrimSuffix(popValue.Value, ","))
	}
	sq.Push(newLine, ") ")
	if qb.Engine == "" {
		panic(errors.New("field Engine can not be empty string"))
	}
	if qb.Engine == "" {
		panic(errors.New("field Engine can not be empty string"))
	}
	if qb.Charset == "" {
		panic(errors.New("field Charset can not be empty string"))
	}
	if qb.Collate == "" {
		panic(errors.New("field Collate can not be empty string"))
	}
	sq.Push("ENGINE=", qb.Engine.String(), " CHARSET=", qb.Charset.String(), " COLLATE=", qb.Collate.String())
	sq.Push(";")
	return sq.Join("")
}
type MigrateField struct {
	name string
	size int
	fieldType string
	unsigned bool
	null bool
	autoIncrement bool
	characterSet string
	collate string
	defaultValue migrateDefaultValue
	references struct{
		valid bool
		otherTableName string
		otherTableField string
	}
	extra []string
	commit string
	raw string
}
func (mi MigrateField) Type(columnType string, size int) MigrateField {
	mi.size = size
	mi.fieldType = columnType
	return mi
}
func (mi MigrateField) Int(size int) MigrateField {
	mi.size = size
	mi.fieldType = "int"
	return mi
}
func (mi MigrateField) Tinyint(size int) MigrateField {
	mi.size = size
	mi.fieldType = "tinyint"
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

func (mi MigrateField) CharacterSet (kind string) MigrateField {
	mi.characterSet = kind
	return mi
}
func (mi MigrateField) Collate(kind string)  MigrateField{
	mi.collate = kind
	return mi
}
type migrateDefaultValue struct {
	raw string
}
func (mi MigrateField) DefaultCurrentTimeStamp() MigrateField {
	mi.defaultValue = migrateDefaultValue{
		raw: "CURRENT_TIMESTAMP",
	}
	return mi
}
func (mi MigrateField) DefaultString(s string) MigrateField {
	mi.defaultValue = migrateDefaultValue{
		raw: `'` + s + `'`,
	}
	return mi
}
func (mi MigrateField) DefaultInt(i int) MigrateField {
	mi.defaultValue = migrateDefaultValue{
		raw: `'` + strconv.Itoa(i) + `'`,
	}
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
func (mi Migrate) Exec(sql string, values... interface{}) {
	_, err := mi.db.Core.DB.Exec(sql, values...)
	mi.CheckError(err, sql)
}
func (mi Migrate) CheckError(err error, sql string) {
	if err != nil {
		log.Print(sql)
		panic(err)
	}
}
func (mi Migrate) CreateTable(qb CreateTableQB) {
	sql := qb.ToSQL()
	_, err := mi.db.Core.DB.Exec(sql) ; mi.CheckError(err, sql)
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
func (mi MigrateField) References(otherTableName string, otherTableField string) MigrateField {
	mi.references.valid = true
	mi.references.otherTableName = otherTableName
	mi.references.otherTableField = otherTableField
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
func (mi Migrate) FieldRaw(raw string) MigrateField {
	if strings.HasSuffix(raw, ",") {
		raw = strings.TrimSuffix(raw, ",")
	}
	return MigrateField{raw: raw}
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
		DefaultCurrentTimeStamp()
}
func (mi Migrate) UpdatedAtTimestamp() MigrateField {
	return mi.Field("updated_at").
		Timestamp().
		DefaultCurrentTimeStamp().
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
