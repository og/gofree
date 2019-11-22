package f

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
	l "github.com/og/x/log"
	"reflect"
)

type Database struct {
	Core *sqlx.DB
	onlyReadDataSourceName DataSourceName
}
func (database Database) GetDataSourceName () DataSourceName {
	return database.onlyReadDataSourceName
}
func NewDatabase(dataSourceName DataSourceName) (database Database) {
	db, err := sqlx.Connect(dataSourceName.DriverName, dataSourceName.GetString())
	ge.Check(err)
	database = Database{
		Core: db,
	}
	database.onlyReadDataSourceName = dataSourceName
	return
}

func (db *Database) OneQB(modelPtr Model, has *bool, qb QB) {
	db.coreOneQB(txOrDB{ UseTx: false,}, modelPtr, has, qb)
	return
}
func (db *Database) TxOneQB(tx *sqlx.Tx, modelPtr Model, has *bool, qb QB) {
	db.coreOneQB(txOrDB{ UseTx: true,Tx: tx,}, modelPtr,has, qb)
	return
}
type txOrDB struct {
	UseTx bool
	Tx *sqlx.Tx
}
func (db *Database) coreOneQB(txDB txOrDB, modelPtr Model, has *bool, qb QB) {
	query, values := qb.BindModel(modelPtr).GetSelect()
	var row *sqlx.Row
	if txDB.UseTx {
		row = txDB.Tx.QueryRowx(query, values...)
	} else {
		row = db.Core.QueryRowx(query, values...)
	}
	err := row.StructScan(modelPtr)
	if err == sql.ErrNoRows { *has = false ; return}
	if err != nil {panic(err)}
	*has = true
	return
}



func (db *Database) OneID(modelPtr Model, has *bool, id interface{}) {
	db.OneQB(modelPtr, has, QB{
		Where:And("id", id),
	})
	return
}

func (db *Database) TxOneID(tx *sqlx.Tx, modelPtr Model, has *bool, id interface{}) {
	db.TxOneQB(tx, modelPtr, has, QB{
		Where:And("id", id),
	})
	return
}
func (db *Database) CountQB(modelPtr Model, qb QB) (count int) {
	return db.coreCountQB(txOrDB{UseTx:false,}, modelPtr, qb)
}
func (db *Database) TxCountQB(tx *sqlx.Tx, modelPtr Model, qb QB) (count int) {
	return db.coreCountQB(txOrDB{UseTx:true, Tx: tx}, modelPtr, qb)
}
func (db *Database) coreCountQB(txDB txOrDB, modelPtr Model, qb QB) (count int) {
	qb.Count = true
	query, values := qb.BindModel(modelPtr).GetSelect()
	var row *sqlx.Row
	if txDB.UseTx {
		row = txDB.Tx.QueryRowx(query, values...)
	} else {
		row = db.Core.QueryRowx(query, values...)
	}
	err := row.Scan(&count)
	if err != nil {panic(err)}
	return
}

func (db *Database) ListQB(modelListPtr interface{}, qb QB) {
	db.coreListQB(txOrDB{UseTx: false}, modelListPtr, qb)
}
func (db *Database) TxListQB(tx *sqlx.Tx, modelListPtr interface{}, qb QB) {
	db.coreListQB(txOrDB{UseTx: false, Tx: tx}, modelListPtr, qb)
}
func (db *Database) coreListQB(txDB txOrDB, modelListPtr interface{}, qb QB) {
	reflectItemValue := reflect.MakeSlice(reflect.TypeOf(modelListPtr).Elem(), 1,1).Index(0)
	model := reflectItemValue.Interface().(Model)
	query, values := qb.BindModel(model).GetSelect()
	if qb.Table == "" {
		tableName := reflectItemValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
	}
	if txDB.UseTx {
		err := txDB.Tx.Select(modelListPtr, query, values...)
		ge.Check(err)
	} else {
		err := db.Core.Select(modelListPtr, query, values...)
		ge.Check(err)
	}
	return
}

func (db *Database) Create(modelPtr interface{}) {
	value := reflect.ValueOf(modelPtr)
	typeValue := reflect.TypeOf(modelPtr)
	fieldLen := value.NumField()
	for i:=0;i<fieldLen;i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		l.V(itemType.Name)
		l.V(itemType.Tag.Get("db"), item.Interface())
	}

	insertData := Map{}
	QB{
		Insert: insertData,
	}.BindModel(modelPtr.(Model)).GetInsert()

}