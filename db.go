package f

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
	"github.com/pkg/errors"
	"reflect"
	"time"
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
	query, values := qb.BindModel(reflectItemValue.Interface()).GetSelect()
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

func (db *Database) coreCreate(txDB txOrDB, modelPtr interface{}) {
	value := reflect.ValueOf(modelPtr).Elem()
	reflect.ValueOf(modelPtr).MethodByName("BeforeCreate").Call([]reflect.Value{})
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := value.NumField()
	insertData := Map{}
	for i:=0;i<fieldLen;i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName := itemType.Tag.Get("db")
		value := item.Interface()
		insertData[dbName] = value
	}
	createdAtValue := value.FieldByName("CreatedAt")
	if createdAtValue.IsValid() {
		createdAtType, _ := typeValue.FieldByName("CreatedAt")
		insertData[createdAtType.Tag.Get("db")] = time.Now()
		createdAtValue.Set(reflect.ValueOf(time.Now()))
	}
	updatedAtValue := value.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		insertData[updatedAtType.Tag.Get("db")] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
	}
	query, values := QB{
		Insert: insertData,
	}.BindModel(modelPtr).GetInsert()
	var result sql.Result
	if txDB.UseTx {
		newResult, err := txDB.Tx.Exec(query, values...) ; ge.Check(err)
		result = newResult
	} else {
		newResult, err := db.Core.Exec(query, values...) ; ge.Check(err)
		result = newResult
	}

	lastInsertID, err := result.LastInsertId() ; ge.Check(err)
	if  lastInsertID != 0 {
		value.FieldByName("ID").SetInt(lastInsertID)
	}

}

func (db *Database) Create(modelPtr interface{}) {
	db.coreCreate(txOrDB{}, modelPtr)
}
func (db *Database) TxCreate(tx *sqlx.Tx, modelPtr interface{}) {
	db.coreCreate(txOrDB{UseTx: true, Tx: tx}, modelPtr)
}

func (db *Database) coreDeleteQB(txDB txOrDB, modelPtr interface{}, qb QB) {
	if len(qb.Update) == 0 {
		qb.Update = Map{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if txDB.UseTx {
		sqlResult, err := txDB.Tx.Exec(query, values...) ; ge.Check(err)
		result = sqlResult
	} else {
		sqlResult, err := db.Core.Exec(query, values...) ; ge.Check(err)
		result = sqlResult
	}
	_, err := result.LastInsertId() ; ge.Check(err)
}
func (db *Database) DeleteQB(modelPtr interface{}, qb QB) {
	db.coreDeleteQB(txOrDB{}, modelPtr, qb)
}
func (db *Database) TxDeleteQB(tx *sqlx.Tx,modelPtr interface{}, qb QB) {
	db.coreDeleteQB(txOrDB{UseTx: true, Tx: tx}, modelPtr, qb)
}


func (db *Database) coreDelete(txDB txOrDB, modelPtr interface{}) {
	id := reflect.ValueOf(modelPtr).Elem().FieldByName("ID").Interface()
	qb := QB{
		Where: And("id", id),
	}
	if len(qb.Update) == 0 {
		qb.Update = Map{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if txDB.UseTx {
		sqlResult, err := txDB.Tx.Exec(query, values...) ; ge.Check(err)
		result = sqlResult
	} else {
		sqlResult, err := db.Core.Exec(query, values...) ; ge.Check(err)
		result = sqlResult
	}
	_, err := result.LastInsertId() ; ge.Check(err)
}
func (db *Database) Delete(modelPtr interface{}) {
	db.coreDelete(txOrDB{}, modelPtr)
}
func (db *Database) TxDelete(tx *sqlx.Tx,modelPtr interface{}, qb QB) {
	db.coreDelete(txOrDB{UseTx: true, Tx: tx}, modelPtr)
}

func (db *Database) Update(modelPtr interface{}) {
	db.coreUpdate(txOrDB{}, modelPtr)
}
func (db *Database) TxUpdate(tx *sqlx.Tx, modelPtr interface{}) {
	db.coreUpdate(txOrDB{UseTx: true, Tx: tx}, modelPtr)
}
func (db *Database) coreUpdate (txDB txOrDB, modelPtr interface{}) {
	value := reflect.ValueOf(modelPtr).Elem()
	reflect.ValueOf(modelPtr).MethodByName("BeforeCreate").Call([]reflect.Value{})
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := value.NumField()
	updateData := Map{}
	var id interface{}
	for i:=0;i<fieldLen;i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName := itemType.Tag.Get("db")
		value := item.Interface()
		if dbName == "id" {
			id = value
			continue
		}
		updateData[dbName] = value
	}
	updatedAtValue := value.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		updateData[updatedAtType.Tag.Get("db")] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
	}
	if id == nil {
		panic(errors.New("db.Update(modelPtr) model.id is nil"))
	}
	query, values := QB{
		Where: And("id", id),
		Update: updateData,
	}.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if txDB.UseTx {
		newResult, err := txDB.Tx.Exec(query, values...) ; ge.Check(err)
		result = newResult
	} else {
		newResult, err := db.Core.Exec(query, values...) ; ge.Check(err)
		result = newResult
	}
	lastInsertID, err := result.LastInsertId() ; ge.Check(err)
	if  lastInsertID != 0 {
		value.FieldByName("ID").SetInt(lastInsertID)
	}
}
