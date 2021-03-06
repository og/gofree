package f

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
	"reflect"
	"time"
)

type Database struct {
	Core *sqlx.DB
	onlyReadDataSourceName DataSourceName
}
func (database Database) Close() {
	ge.Check(database.Core.Close())
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
func (db *Database) TxOneQB(tx *Tx, modelPtr Model, has *bool, qb QB) {
	db.coreOneQB(txOrDB{ UseTx: true,Tx: tx.core,}, modelPtr,has, qb)
	return
}
type txOrDB struct {
	UseTx bool
	Tx *sqlx.Tx
}
func (db *Database) coreOneQB(txDB txOrDB, modelPtr Model, has *bool, qb QB) {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.OneID() or db.OneQB()  arg `modelPtr` must be a ptr")
	}
	scanModelMakeSQLSelect(reflect.ValueOf(modelPtr).Elem().Type(), &qb)
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

func (db *Database) TxOneID(tx *Tx, modelPtr Model, has *bool, id interface{}) {
	db.TxOneQB(tx, modelPtr, has, QB{
		Where:And("id", id),
	})
	return
}
func (db *Database) CountQB(modelPtr Model, qb QB) (count int) {
	return db.coreCountQB(txOrDB{UseTx:false,}, modelPtr, qb)
}
func (db *Database) TxCountQB(tx *Tx, modelPtr Model, qb QB) (count int) {
	return db.coreCountQB(txOrDB{UseTx:true, Tx: tx.core}, modelPtr, qb)
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
func (db *Database) TxListQB(tx *Tx, modelListPtr []Model, qb QB) {
	db.coreListQB(txOrDB{UseTx: false, Tx: tx.core}, modelListPtr, qb)
}
func (db *Database) coreListQB(txDB txOrDB, modelListPtr interface{}, qb QB) {
	elementType := reflect.TypeOf(modelListPtr).Elem()
	reflectItemValue := reflect.MakeSlice(elementType, 1,1).Index(0)
	modelInterface := reflectItemValue.Interface()
	scanModelMakeSQLSelect(reflectItemValue.Type(), &qb)
	query, values := qb.BindModel(modelInterface).GetSelect()
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

func (db *Database) coreCreate(txDB txOrDB, modelPtr Model) {
	value, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.OneID() or db.OneQB()  arg `modelPtr` must be a ptr")
	}
	reflect.ValueOf(modelPtr).MethodByName("BeforeCreate").Call([]reflect.Value{})
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := value.NumField()
	insertData := Map{}
	for i:=0;i<fieldLen;i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName := itemType.Tag.Get("db")
		if dbName == "" {
			continue
		}
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

func (db *Database) Create(modelPtr Model) {
	db.coreCreate(txOrDB{}, modelPtr)
}
func (db *Database) TxCreate(tx *Tx, modelPtr Model) {
	db.coreCreate(txOrDB{UseTx: true, Tx: tx.core}, modelPtr)
}

func (db *Database) coreDeleteQB(txDB txOrDB, modelPtr Model, qb QB) {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.DeleteQB() or db.TxDeleteQB()  arg `modelPtr` must be a ptr, eg: db.DeleteQB(&user, qb) db.TxDeleteQB(tx, &user, qb) ")
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
func (db *Database) DeleteQB(modelPtr Model, qb QB) {
	db.coreDeleteQB(txOrDB{}, modelPtr, qb)
}
func (db *Database) TxDeleteQB(tx *Tx,modelPtr Model, qb QB) {
	db.coreDeleteQB(txOrDB{UseTx: true, Tx: tx.core}, modelPtr, qb)
}


func (db *Database) coreDelete(txDB txOrDB, modelPtr Model) {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.Delete() or db.TxDelete()  arg `modelPtr` must be a ptr, eg: db.Delete(&user) db.TxDelete(tx, &user) ")
	}
	idValue := rValue.FieldByName("ID")
	if idValue.IsZero() {
		panic(errors.New("db.Update(&model) or db.TxUpdate(&model) model.id is zero"))
	}
	id := idValue.Interface()


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
func (db *Database) Delete(modelPtr Model) {
	db.coreDelete(txOrDB{}, modelPtr)
}
func (db *Database) TxDelete(tx *Tx,modelPtr Model, qb QB) {
	db.coreDelete(txOrDB{UseTx: true, Tx: tx.core}, modelPtr)
}

func (db *Database) Update(modelPtr Model) {
	db.coreUpdate(txOrDB{}, modelPtr)
}
func (db *Database) TxUpdate(tx *Tx, modelPtr Model) {
	db.coreUpdate(txOrDB{UseTx: true, Tx: tx.core}, modelPtr)
}
func getPtrElem(ptr interface{}) (value reflect.Value, isPtr bool) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		isPtr = false
		return
	}
	value = v.Elem()
	isPtr = true
	return
}
func (db *Database) coreUpdate (txDB txOrDB, modelPtr Model) {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.Update() or db.TxUpdate()  arg `modelPtr` must be a ptr, eg: db.Update(&user) db.TxUpdate(tx, &user) ")
	}
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := rValue.NumField()
	updateData := Map{}
	var id interface{}
	for i:=0;i<fieldLen;i++{
		item := rValue.Field(i)
		itemType := typeValue.Field(i)
		dbName := itemType.Tag.Get("db")
		value := item.Interface()
		if dbName == "id" {
			if item.IsZero() {
				panic(errors.New("db.Update(&model) or db.TxUpdate(&model) model.id is zero"))
			}
			id = value
			continue
		}
		updateData[dbName] = value
	}
	updatedAtValue := rValue.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		updateData[updatedAtType.Tag.Get("db")] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
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
		rValue.FieldByName("ID").SetInt(lastInsertID)
	}
}

func (db Database) Tx() *Tx {
	tx, err := db.Core.Beginx() ; ge.Check(err)
	return newTx(tx)
}
// 配置项
type ReadQB struct {
	Limit int
	ListPtr interface{}
	QB QB
	Read func()
	Notes NotesReadQBType
}
type NotesReadQBType struct {
	Message string
}
var NotesReadQB = NotesReadQBType {
	Message: `
请确保基于 f.QB 的查询结果不可变，或者基于时间倒叙查询。（从最早的数据查询到最新的数据）
目的是为了避免查询数据过程中有新的数据插入，此时 ReadQB 内部通过 offset 查询时会不准确。 
`,
}
// must use asc by order created at or query results are immutable
func (db Database) ReadQB(conf ReadQB) {
	qb := conf.QB
	qb.Limit = conf.Limit
	qb.Offset = 0
	listPtr := conf.ListPtr
	for {
		list := reflect.ValueOf(listPtr).Elem()
		list.SetLen(0)
		db.ListQB(listPtr, qb)
		if list.Len() == 0 {
			break
		}
		conf.Read()
		qb.Offset += qb.Limit
	}
}
