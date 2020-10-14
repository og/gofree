package f

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
	"log"
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
func (database Database) DataSourceName () DataSourceName {
	return database.onlyReadDataSourceName
}
func NewDatabase(dataSourceName DataSourceName) (database Database, err error) {
	db, err := sqlx.Connect(dataSourceName.DriverName, dataSourceName.GetString())
	if err != nil {
		return database, err
	}
	database = Database{Core: db,}
	database.onlyReadDataSourceName = dataSourceName
	return database, nil
}
type SqlOpt struct {
	Tx *Tx
	Context context.Context
}
func (opt SqlOpt) ctxOrTODO() context.Context {
	if opt.Context == nil {
		return context.TODO()
	} else {
		return opt.Context
	}
}
func (database *Database) OneQB(modelPtr Model, has *bool, qb QB) {
	ge.Check(database.CoreOneQB(SqlOpt{}, modelPtr, has, qb))
}
func (database *Database) CoreOneQB(opt SqlOpt, modelPtr Model, has *bool, qb QB) error{
	qb.Limit = 1
	scanModelMakeSQLSelect(reflect.ValueOf(modelPtr).Elem().Type(), &qb)
	query, values := qb.BindModel(modelPtr).GetSelect()
	var row *sqlx.Row
	if opt.Tx != nil  {
		row = opt.Tx.core.QueryRowxContext(opt.ctxOrTODO(), query, values...)
	} else {
		row = database.Core.QueryRowxContext(opt.ctxOrTODO(),query, values...)
	}
	err := row.StructScan(modelPtr)
	if err == sql.ErrNoRows { *has = false ; return nil}
	if err != nil {
		return err
	}
	*has = true
	return nil
}



func (database *Database) OneID(modelPtr Model, has *bool, id interface{}) {
	database.OneQB(modelPtr, has, QB{
		Where:And("id", id),
	})
}

func (database *Database) CountQB(modelPtr Model, qb QB) (count int)  {
	return database.coreCountQB(SqlOpt{}, modelPtr, qb)
}
func (database *Database) CountQBContext(ctx context.Context, modelPtr Model, qb QB) (count int)  {
	return database.coreCountQB(SqlOpt{Context: ctx}, modelPtr, qb)
}
func (database *Database) TxCountQB(tx *Tx, modelPtr Model, qb QB) (count int) {
	return database.coreCountQB(SqlOpt{Tx: tx}, modelPtr, qb)
}
func (database *Database) TxCountQBContext(ctx context.Context, tx *Tx, modelPtr Model, qb QB) (count int) {
	return database.coreCountQB(SqlOpt{Tx: tx, Context: ctx}, modelPtr, qb)
}
func (database *Database) coreCountQB(opt SqlOpt, modelPtr Model, qb QB) (count int)  {
	qb.Count = true
	query, values := qb.BindModel(modelPtr).GetSelect()
	var row *sqlx.Row
	if opt.Tx != nil {
		row = opt.Tx.core.QueryRowxContext(opt.ctxOrTODO(),query, values...)
	} else {
		row = database.Core.QueryRowxContext(opt.ctxOrTODO(),query, values...)
	}
	err := row.Scan(&count)
	if err != nil {panic(err)}
	return
}

func (database *Database) ListQB(modelListPtr interface{}, qb QB) {
	database.coreListQB(SqlOpt{}, modelListPtr, qb)
}
func (database *Database) ListQBContext(ctx context.Context, modelListPtr interface{}, qb QB) {
	database.coreListQB(SqlOpt{Context: ctx}, modelListPtr, qb)
}
func (database *Database) TxListQB(tx *Tx, modelListPtr []Model, qb QB) {
	database.coreListQB(SqlOpt{Tx: tx}, modelListPtr, qb)
}
func (database *Database) TxListQBContext(ctx context.Context, tx *Tx, modelListPtr []Model, qb QB) {
	database.coreListQB(SqlOpt{Tx: tx, Context:ctx}, modelListPtr, qb)
}
func (database *Database) coreListQB(opt SqlOpt, modelListPtr interface{}, qb QB) {
	elementType := reflect.TypeOf(modelListPtr).Elem()
	reflectItemValue := reflect.MakeSlice(elementType, 1,1).Index(0)
	scanModelMakeSQLSelect(reflectItemValue.Type(), &qb)
	query, values := qb.BindModel(reflectItemValue.Addr().Interface().(Model)).GetSelect()
	if qb.Table == "" {
		tableName := reflectItemValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
		qb.Table = tableName
	}
	if opt.Tx != nil {
		err := opt.Tx.core.SelectContext(opt.ctxOrTODO(),modelListPtr, query, values...)
		ge.Check(err)
	} else {
		err := database.Core.SelectContext(opt.ctxOrTODO(),modelListPtr, query, values...)
		ge.Check(err)
	}
	return
}

func (database *Database) CoreCreate(opt SqlOpt, modelPtr Model) error {
	value, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.OneID() or db.OneQB()  arg `modelPtr` must be a ptr")
	}
	modelPtr.BeforeCreate()
	typeValue := reflect.TypeOf(modelPtr).Elem()
	insertData := map[Column]interface{}{}
	for i:=0;i<value.NumField();i++{
		item := value.Field(i)
		itemType := typeValue.Field(i)
		dbName, hasDBName := itemType.Tag.Lookup("db")
		if !hasDBName {
			continue
		}
		if dbName == "" {
			log.Print(`Maybe you forget set db:"name"` + itemType.Name)
			continue
		}
		insertData[Column(dbName)] = item.Interface()
	}
	createdAtValue := value.FieldByName("CreatedAt")
	if createdAtValue.IsValid() {
		createdAtType, _ := typeValue.FieldByName("CreatedAt")
		insertData[Column(createdAtType.Tag.Get("db"))] = time.Now()
		createdAtValue.Set(reflect.ValueOf(time.Now()))
	}
	updatedAtValue := value.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		insertData[Column(updatedAtType.Tag.Get("db"))] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
	}
	query, values := QB{
		Insert: insertData,
	}.BindModel(modelPtr).GetInsert()
	var result sql.Result
	if opt.Tx != nil {
		newResult, err := opt.Tx.core.ExecContext(opt.ctxOrTODO(),query, values...) ; if err != nil {return err}
		result = newResult
	} else {
		newResult, err := database.Core.ExecContext(opt.ctxOrTODO(),query, values...) ; if err != nil {return err}
		result = newResult
	}
	lastInsertID, err := result.LastInsertId() ; ge.Check(err)
	if  lastInsertID != 0 {
		value.FieldByName("ID").SetInt(lastInsertID)
	}
	return nil
}

func (database *Database) Create(modelPtr Model) {
	err := database.CoreCreate(SqlOpt{}, modelPtr)
	if err != nil { panic(err) }
}

func (database *Database) coreDeleteQB(opt SqlOpt, modelPtr Model, qb QB) {
	_, isPtr := getPtrElem(modelPtr)
	if !isPtr {
		panic("db.DeleteQB() or db.TxDeleteQB()  arg `modelPtr` must be a ptr, eg: db.DeleteQB(&user, qb) db.TxDeleteQB(tx, &user, qb) ")
	}
	if len(qb.Update) == 0 {
		qb.Update = map[Column]interface{}{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if opt.Tx != nil {
		sqlResult, err := opt.Tx.core.ExecContext(opt.ctxOrTODO(),query, values...) ; ge.Check(err)
		result = sqlResult
	} else {
		sqlResult, err := database.Core.ExecContext(opt.ctxOrTODO(),query, values...) ; ge.Check(err)
		result = sqlResult
	}
	_, err := result.LastInsertId() ; ge.Check(err)
}
func (database *Database) DeleteQB(modelPtr Model, qb QB) {
	database.coreDeleteQB(SqlOpt{}, modelPtr, qb)
}
func (database *Database) DeleteQBContext(ctx context.Context, modelPtr Model, qb QB) {
	database.coreDeleteQB(SqlOpt{Context: ctx}, modelPtr, qb)
}
func (database *Database) TxDeleteQB(tx *Tx,modelPtr Model, qb QB) {
	database.coreDeleteQB(SqlOpt{Tx: tx}, modelPtr, qb)
}
func (database *Database) TxDeleteQBContext(ctx context.Context, tx *Tx,modelPtr Model, qb QB) {
	database.coreDeleteQB(SqlOpt{Tx: tx, Context:ctx}, modelPtr, qb)
}


func (database *Database) coreDelete(opt SqlOpt, modelPtr Model) {
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
		qb.Update = map[Column]interface{}{}
	}
	qb.Update["deleted_at"] = time.Now()
	query, values := qb.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if opt.Tx != nil {
		sqlResult, err := opt.Tx.core.ExecContext(opt.ctxOrTODO(),query, values...) ; ge.Check(err)
		result = sqlResult
	} else {
		sqlResult, err := database.Core.ExecContext(opt.ctxOrTODO(),query, values...) ; ge.Check(err)
		result = sqlResult
	}
	_, err := result.LastInsertId() ; ge.Check(err)
}
func (database *Database) Delete(modelPtr Model) {
	database.coreDelete(SqlOpt{}, modelPtr)
}
func (database *Database) DeleteContext(ctx context.Context, modelPtr Model) {
	database.coreDelete(SqlOpt{Context: ctx,}, modelPtr)
}
func (database *Database) TxDelete(tx *Tx,modelPtr Model, qb QB) {
	database.coreDelete(SqlOpt{Tx: tx}, modelPtr)
}
func (database *Database) TxDeleteContext(ctx context.Context, tx *Tx,modelPtr Model, qb QB) {
	database.coreDelete(SqlOpt{Tx: tx, Context:ctx}, modelPtr)
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

func (database *Database) baseUpdate (opt SqlOpt, modelPtr Model, useUpdateData bool, updateData Data) error {
	rValue, isPtr := getPtrElem(modelPtr)
	if !isPtr { panic("db.Update() or db.TxUpdate()  arg `modelPtr` must be a ptr, eg: db.Update(&user) db.TxUpdate(tx, &user) ") }
	typeValue := reflect.TypeOf(modelPtr).Elem()
	fieldLen := rValue.NumField()
	var id interface{}
	var findID bool
	for i:=0;i<fieldLen;i++{
		item := rValue.Field(i)
		itemType := typeValue.Field(i)
		dbName, hasDBName := itemType.Tag.Lookup("db")
		if !hasDBName {
			continue
		}
		if dbName == "" {
			log.Print(`Maybe you forget set db:"name"` + itemType.Name)
			continue
		}
		value := item.Interface()
		if dbName == "id" {
			if item.IsZero() {
				panic(errors.New("update model.ID is zero"))
			}
			findID = true
			id = value
			continue
		}
		if useUpdateData {
			rValue.Set(reflect.ValueOf(updateData[Column(dbName)]))
		} else {
			updateData[Column(dbName)] = value
		}
	}
	if !findID {
		panic(errors.New("update not found id"))
	}
	updatedAtValue := rValue.FieldByName("UpdatedAt")
	if updatedAtValue.IsValid() {
		updatedAtType, _ := typeValue.FieldByName("UpdatedAt")
		updateData[Column(updatedAtType.Tag.Get("db"))] = time.Now()
		updatedAtValue.Set(reflect.ValueOf(time.Now()))
	}

	query, values := QB{
		Where: And("id", id),
		Update: updateData,
	}.BindModel(modelPtr).GetUpdate()
	var result sql.Result
	if opt.Tx != nil {
		newResult, err := opt.Tx.core.ExecContext(opt.ctxOrTODO(),query, values...) ; if err != nil {return err}
		result = newResult
	} else {
		newResult, err := database.Core.ExecContext(opt.ctxOrTODO(),query, values...) ; if err != nil {return err}
		result = newResult
	}
	lastInsertID, err := result.LastInsertId(); if err != nil {return err}
	if  lastInsertID != 0 {
		rValue.FieldByName("ID").SetInt(lastInsertID)
	}
	return nil
}
func (database *Database) UpdateData (modelPtr Model, data Data){
	err := database.CoreUpdateData(SqlOpt{}, modelPtr, data)
	if err != nil {panic(err)}
}
func (database *Database) CoreUpdateData (opt SqlOpt, modelPtr Model, data Data) error {
	return database.baseUpdate(opt, modelPtr, true, data)
}
func (database *Database) Update(modelPtr Model) {
	err := database.CoreUpdate(SqlOpt{}, modelPtr)
	if err != nil {panic(err)}
}
func (database *Database) CoreUpdate (opt SqlOpt, modelPtr Model) error {
	return database.baseUpdate(opt, modelPtr, false, nil)
}

func (database Database) Tx(transaction func(tx *Tx) error) error {
	return database.CoreTx(context.Background(), nil, transaction)
}
func (database Database) CoreTx(ctx context.Context, options *sql.TxOptions, transaction func(tx *Tx) error ) error {
	dbTx, err := database.Core.BeginTxx(ctx, options)
	if err != nil {return err}
	tx := newTx(dbTx)
	var txErr error
	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
			panic(r)
		} else {
			if txErr == nil {
				tx.Commit()
			}
		}
	}()
	txErr = transaction(tx)
	if txErr != nil {
		tx.Rollback()
		return txErr
	}
	return nil
}